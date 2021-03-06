---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_database"
sidebar_current: "docs-azurerm-resource-database-sql-database"
description: |-
  Create a SQL Database.
---

# azurerm_sql_database

Allows you to manage an Azure SQL Database

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_sql_server" "test" {
    name = "mysqlserver"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    version = "12.0"
    administrator_login = "4dm1n157r470r"
    administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_database" "test" {
  name                = "mysqldatabase"
  resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    server_name = "${azurerm_sql_server.test.name}"

  tags {
    environment = "production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database.

* `resource_group_name` - (Required) The name of the resource group in which to create the database.  This must be the same as Database Server resource group currently.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `server_name` - (Required) The name of the SQL Server on which to create the database.

* `create_mode` - (Optional) Specifies the type of database to create. Defaults to `Default`. See below for the accepted values/

* `import` - (Optional) A Database Import block as documented below. `create_mode` must be set to `Default`.

* `source_database_id` - (Optional) The URI of the source database if `create_mode` value is not `Default`.

* `restore_point_in_time` - (Optional) The point in time for the restore. Only applies if `create_mode` is `PointInTimeRestore` e.g. 2013-11-08T22:00:40Z

* `edition` - (Optional) The edition of the database to be created. Applies only if `create_mode` is `Default`. Valid values are: `Basic`, `Standard`, `Premium`, or `DataWarehouse`. Please see [Azure SQL Database Service Tiers](https://azure.microsoft.com/en-gb/documentation/articles/sql-database-service-tiers/).

* `collation` - (Optional) The name of the collation. Applies only if `create_mode` is `Default`.  Azure default is `SQL_LATIN1_GENERAL_CP1_CI_AS`. Changing this forces a new resource to be created.

* `max_size_bytes` - (Optional) The maximum size that the database can grow to. Applies only if `create_mode` is `Default`.  Please see [Azure SQL Database Service Tiers](https://azure.microsoft.com/en-gb/documentation/articles/sql-database-service-tiers/).

* `requested_service_objective_id` - (Optional) Use `requested_service_objective_id` or `requested_service_objective_name` to set the performance level for the database.
 Valid values are: `S0`, `S1`, `S2`, `S3`, `P1`, `P2`, `P4`, `P6`, `P11` and `ElasticPool`.  Please see [Azure SQL Database Service Tiers](https://azure.microsoft.com/en-gb/documentation/articles/sql-database-service-tiers/).

* `requested_service_objective_name` - (Optional) Use `requested_service_objective_name` or `requested_service_objective_id` to set the performance level for the database.  Please see [Azure SQL Database Service Tiers](https://azure.microsoft.com/en-gb/documentation/articles/sql-database-service-tiers/).

* `source_database_deletion_date` - (Optional) The deletion date time of the source database. Only applies to deleted databases where `create_mode` is `PointInTimeRestore`.

* `elastic_pool_name` - (Optional) The name of the elastic database pool.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`import` supports the following:

* `storage_uri` - (Required) Specifies the blob URI of the .bacpac file.
* `storage_key` - (Required) Specifies the access key for the storage account.
* `storage_key_type` - (Required) Specifies the type of access key for the storage account. Valid values are `StorageAccessKey` or `SharedAccessKey`.
* `administrator_login` - (Required) Specifies the name of the SQL administrator.
* `administrator_login_password` - (Required) Specifies the password of the SQL administrator.
* `authentication_type` - (Required) Specifies the type of authentication used to access the server. Valid values are `SQL` or `ADPassword`.
* `operation_mode` - (Optional) Specifies the type of import operation being performed. The only allowable value is `Import`.

## Attributes Reference

The following attributes are exported:

* `id` - The SQL Database ID.
* `creation_date` - The creation date of the SQL Database.
* `default_secondary_location` - The default secondary location of the SQL Database.

## Import

SQL Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_database.database1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/databases/database1
```
