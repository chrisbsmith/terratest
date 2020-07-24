package aws

import (
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/chrisbsmith/terratest/modules/testing"
	_ "github.com/go-sql-driver/mysql"
)

// GetAddressOfRdsInstance gets the address of the given RDS Instance in the given region.
func GetAddressOfRdsInstance(t testing.TestingT, dbInstanceID string, awsRegion string) string {
	address, err := GetAddressOfRdsInstanceE(t, dbInstanceID, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return address
}

// GetAddressOfRdsInstanceE gets the address of the given RDS Instance in the given region.
func GetAddressOfRdsInstanceE(t testing.TestingT, dbInstanceID string, awsRegion string) (string, error) {
	dbInstance, err := GetRdsInstanceDetailsE(t, dbInstanceID, awsRegion)
	if err != nil {
		return "", err
	}

	return aws.StringValue(dbInstance.Endpoint.Address), nil
}

// GetPortOfRdsInstance gets the address of the given RDS Instance in the given region.
func GetPortOfRdsInstance(t testing.TestingT, dbInstanceID string, awsRegion string) int64 {
	port, err := GetPortOfRdsInstanceE(t, dbInstanceID, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return port
}

// GetPortOfRdsInstanceE gets the address of the given RDS Instance in the given region.
func GetPortOfRdsInstanceE(t testing.TestingT, dbInstanceID string, awsRegion string) (int64, error) {
	dbInstance, err := GetRdsInstanceDetailsE(t, dbInstanceID, awsRegion)
	if err != nil {
		return -1, err
	}

	return *dbInstance.Endpoint.Port, nil
}

// GetWhetherSchemaExistsInRdsMySqlInstance checks whether the specified schema/table name exists in the RDS instance
func GetWhetherSchemaExistsInRdsMySqlInstance(t testing.TestingT, dbUrl string, dbPort int64, dbUsername string, dbPassword string, expectedSchemaName string) bool {
	output, err := GetWhetherSchemaExistsInRdsMySqlInstanceE(t, dbUrl, dbPort, dbUsername, dbPassword, expectedSchemaName)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

// GetWhetherSchemaExistsInRdsMySqlInstanceE checks whether the specified schema/table name exists in the RDS instance
func GetWhetherSchemaExistsInRdsMySqlInstanceE(t testing.TestingT, dbUrl string, dbPort int64, dbUsername string, dbPassword string, expectedSchemaName string) (bool, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/", dbUsername, dbPassword, dbUrl, dbPort)
	db, connErr := sql.Open("mysql", connectionString)
	if connErr != nil {
		return false, connErr
	}
	defer db.Close()
	var (
		schemaName string
	)
	sqlStatement := "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME=?;"
	row := db.QueryRow(sqlStatement, expectedSchemaName)
	scanErr := row.Scan(&schemaName)
	if scanErr != nil {
		return false, scanErr
	}
	return true, nil
}

// GetParameterValueForParameterOfRdsInstance gets the value of the parameter name specified for the RDS instance in the given region.
func GetParameterValueForParameterOfRdsInstance(t testing.TestingT, parameterName string, dbInstanceID string, awsRegion string) string {
	parameterValue, err := GetParameterValueForParameterOfRdsInstanceE(t, parameterName, dbInstanceID, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return parameterValue
}

// GetParameterValueForParameterOfRdsInstanceE gets the value of the parameter name specified for the RDS instance in the given region.
func GetParameterValueForParameterOfRdsInstanceE(t testing.TestingT, parameterName string, dbInstanceID string, awsRegion string) (string, error) {
	output := GetAllParametersOfRdsInstance(t, dbInstanceID, awsRegion)
	for _, parameter := range output {
		if aws.StringValue(parameter.ParameterName) == parameterName {
			return aws.StringValue(parameter.ParameterValue), nil
		}
	}
	return "", ParameterForDbInstanceNotFound{ParameterName: parameterName, DbInstanceID: dbInstanceID, AwsRegion: awsRegion}
}

// GetOptionSettingForOfRdsInstance gets the value of the option name in the option group specified for the RDS instance in the given region.
func GetOptionSettingForOfRdsInstance(t testing.TestingT, optionName string, optionSettingName string, dbInstanceID, awsRegion string) string {
	optionValue, err := GetOptionSettingForOfRdsInstanceE(t, optionName, optionSettingName, dbInstanceID, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return optionValue
}

// GetOptionSettingForOfRdsInstanceE gets the value of the option name in the option group specified for the RDS instance in the given region.
func GetOptionSettingForOfRdsInstanceE(t testing.TestingT, optionName string, optionSettingName string, dbInstanceID, awsRegion string) (string, error) {
	optionGroupName := GetOptionGroupNameOfRdsInstance(t, dbInstanceID, awsRegion)
	options := GetOptionsOfOptionGroup(t, optionGroupName, awsRegion)
	for _, option := range options {
		if aws.StringValue(option.OptionName) == optionName {
			for _, optionSetting := range option.OptionSettings {
				if aws.StringValue(optionSetting.Name) == optionSettingName {
					return aws.StringValue(optionSetting.Value), nil
				}
			}
		}
	}
	return "", OptionGroupOptionSettingForDbInstanceNotFound{OptionName: optionName, OptionSettingName: optionSettingName, DbInstanceID: dbInstanceID, AwsRegion: awsRegion}
}

// GetOptionGroupNameOfRdsInstance gets the name of the option group associated with the RDS instance
func GetOptionGroupNameOfRdsInstance(t testing.TestingT, dbInstanceID string, awsRegion string) string {
	dbInstance, err := GetOptionGroupNameOfRdsInstanceE(t, dbInstanceID, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return dbInstance
}

// GetOptionGroupNameOfRdsInstanceE gets the name of the option group associated with the RDS instance
func GetOptionGroupNameOfRdsInstanceE(t testing.TestingT, dbInstanceID string, awsRegion string) (string, error) {
	dbInstance, err := GetRdsInstanceDetailsE(t, dbInstanceID, awsRegion)
	if err != nil {
		return "", err
	}
	return aws.StringValue(dbInstance.OptionGroupMemberships[0].OptionGroupName), nil
}

// GetOptionsOfOptionGroup gets the options of the option group specified
func GetOptionsOfOptionGroup(t testing.TestingT, optionGroupName string, awsRegion string) []*rds.Option {
	output, err := GetOptionsOfOptionGroupE(t, optionGroupName, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

// GetOptionsOfOptionGroupE gets the options of the option group specified
func GetOptionsOfOptionGroupE(t testing.TestingT, optionGroupName string, awsRegion string) ([]*rds.Option, error) {
	rdsClient := NewRdsClient(t, awsRegion)
	input := rds.DescribeOptionGroupsInput{OptionGroupName: aws.String(optionGroupName)}
	output, err := rdsClient.DescribeOptionGroups(&input)
	if err != nil {
		return []*rds.Option{}, err
	}
	return output.OptionGroupsList[0].Options, nil
}

// GetAllParametersOfRdsInstance gets all the parameters defined in the parameter group for the RDS instance in the given region.
func GetAllParametersOfRdsInstance(t testing.TestingT, dbInstanceID string, awsRegion string) []*rds.Parameter {
	parameters, err := GetAllParametersOfRdsInstanceE(t, dbInstanceID, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return parameters
}

// GetAllParametersOfRdsInstanceE gets all the parameters defined in the parameter group for the RDS instance in the given region.
func GetAllParametersOfRdsInstanceE(t testing.TestingT, dbInstanceID string, awsRegion string) ([]*rds.Parameter, error) {
	dbInstance, dbInstanceErr := GetRdsInstanceDetailsE(t, dbInstanceID, awsRegion)
	if dbInstanceErr != nil {
		return []*rds.Parameter{}, dbInstanceErr
	}
	parameterGroupName := aws.StringValue(dbInstance.DBParameterGroups[0].DBParameterGroupName)

	rdsClient := NewRdsClient(t, awsRegion)
	input := rds.DescribeDBParametersInput{DBParameterGroupName: aws.String(parameterGroupName)}
	output, err := rdsClient.DescribeDBParameters(&input)

	if err != nil {
		return []*rds.Parameter{}, err
	}
	return output.Parameters, nil
}

// GetRdsInstanceDetailsE gets the details of a single DB instance whose identifier is passed.
func GetRdsInstanceDetailsE(t testing.TestingT, dbInstanceID string, awsRegion string) (*rds.DBInstance, error) {
	rdsClient := NewRdsClient(t, awsRegion)
	input := rds.DescribeDBInstancesInput{DBInstanceIdentifier: aws.String(dbInstanceID)}
	output, err := rdsClient.DescribeDBInstances(&input)
	if err != nil {
		return nil, err
	}
	return output.DBInstances[0], nil
}

// NewRdsClient creates an RDS client.
func NewRdsClient(t testing.TestingT, region string) *rds.RDS {
	client, err := NewRdsClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewRdsClientE creates an RDS client.
func NewRdsClientE(t testing.TestingT, region string) (*rds.RDS, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return rds.New(sess), nil
}

// ParameterForDbInstanceNotFound is an error that occurs when the parameter group specified is not found for the DB instance
type ParameterForDbInstanceNotFound struct {
	ParameterName string
	DbInstanceID  string
	AwsRegion     string
}

func (err ParameterForDbInstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a parameter %s in parameter group of database %s in %s", err.ParameterName, err.DbInstanceID, err.AwsRegion)
}

// OptionGroupOptionSettingForDbInstanceNotFound is an error that occurs when the option setting specified is not found in the option group of the DB instance
type OptionGroupOptionSettingForDbInstanceNotFound struct {
	OptionName        string
	OptionSettingName string
	DbInstanceID      string
	AwsRegion         string
}

func (err OptionGroupOptionSettingForDbInstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a option setting %s in option name %s of database %s in %s", err.OptionName, err.OptionSettingName, err.DbInstanceID, err.AwsRegion)
}
