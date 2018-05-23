package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultAccessPolicy_basic(t *testing.T) {
	resourceName := "azurerm_key_vault_access_policy.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultAccessPolicy_basic(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_complete(t *testing.T) {
	resourceName := "azurerm_key_vault_access_policy.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultAccessPolicy_complete(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "create"),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.1", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "delete"),
					resource.TestCheckResourceAttr(resourceName, "certificate_permissions.0", "create"),
					resource.TestCheckResourceAttr(resourceName, "certificate_permissions.1", "delete"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_update(t *testing.T) {
	rs := acctest.RandString(6)
	resourceName := "azurerm_key_vault_access_policy.test"
	preConfig := testAccAzureRMKeyVaultAccessPolicy_basic(rs, testLocation())
	postConfig := testAccAzureRMKeyVaultAccessPolicy_update(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// Using the KeyVaultDestroy checker since that ensures the policy is gone.
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists("azurerm_key_vault_access_policy.test"),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.1", "encrypt"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_policyRemoved(t *testing.T) {
	rs := acctest.RandString(6)
	resourceName := "azurerm_key_vault_access_policy.test"
	preConfig := testAccAzureRMKeyVaultAccessPolicy_basic(rs, testLocation())
	postConfig := testAccAzureRMKeyVaultAccessPolicy_policyRemoved(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// Using the KeyVaultDestroy checker since that ensures the policy is gone.
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyMissing("azurerm_key_vault_access_policy.test"),
				),
			},
		},
	})
}

func testCheckAzureRMKeyVaultAccessPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		name := rs.Primary.Attributes["vault_name"]
		resGroup := rs.Primary.Attributes["vault_resource_group"]

		client := testAccProvider.Meta().(*ArmClient).keyVaultClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault %q (resource group: %q) does not exist", name, resGroup)
			}

			return fmt.Errorf("Bad: Get on keyVaultClient: %+v", err)
		}

		objectId := rs.Primary.Attributes["object_id"]
		policy := findKeyVaultAccessPolicy(objectId, resp.Properties.AccessPolicies)

		if policy == nil {
			return fmt.Errorf("Bad: Key Vault Policy %q (resource group: %q, object_id: %s) does not exist", name, resGroup, objectId)
		}

		return nil
	}
}

func testCheckAzureRMKeyVaultAccessPolicyMissing(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).keyVaultClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault %q (resource group: %q) does not exist", name, resGroup)
			}

			return fmt.Errorf("Bad: Get on keyVaultClient: %+v", err)
		}

		objectId := rs.Primary.Attributes["tags.policy_object_id"]

		policy := findKeyVaultAccessPolicy(objectId, resp.Properties.AccessPolicies)

		if policy != nil {
			return fmt.Errorf("Bad: Key Vault Policy %q (resource group: %q, object_id: %s) exists", name, resGroup, objectId)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultAccessPolicy_basic(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
	vault_name                = "${azurerm_key_vault.test.name}"
	vault_resource_group      = "${azurerm_resource_group.test.name}"

	key_permissions = [
      "get"
	]

    secret_permissions = [
      "get",
      "set"
    ]

	tenant_id = "${data.azurerm_client_config.current.tenant_id}"
	object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, rString, location, rString)
}

func testAccAzureRMKeyVaultAccessPolicy_complete(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"


  sku {
    name = "premium"
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
	vault_name                = "${azurerm_key_vault.test.name}"
	vault_resource_group      = "${azurerm_resource_group.test.name}"

	key_permissions = [
	  "create",
      "get"
	]

    secret_permissions = [
      "get",
      "delete"
    ]

	certificate_permissions = [
		"create", 
		"delete"
	]

	application_id 	= "${data.azurerm_client_config.current.service_principal_application_id}"
	tenant_id		= "${data.azurerm_client_config.current.tenant_id}"
	object_id 		= "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, rString, location, rString)
}

func testAccAzureRMKeyVaultAccessPolicy_update(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
	vault_name                = "${azurerm_key_vault.test.name}"
	vault_resource_group      = "${azurerm_resource_group.test.name}"

	key_permissions = [
      "list",
      "encrypt"
	]

    secret_permissions = [
    ]

	tenant_id = "${data.azurerm_client_config.current.tenant_id}"
	object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, rString, location, rString)
}

func testAccAzureRMKeyVaultAccessPolicy_policyRemoved(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  tags {
    environment = "Production"
	policy_object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
  }
}

`, rString, location, rString)
}