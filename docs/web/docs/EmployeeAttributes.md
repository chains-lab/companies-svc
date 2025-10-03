# EmployeeAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DistributorId** | [**uuid.UUID**](uuid.UUID.md) | The unique identifier for the distributor associated with the employee. | 
**Role** | **string** | The role of the employee within the distributor&#39;s organization (e.g., manager, staff). | 
**CreatedAt** | **time.Time** | The timestamp when the employee record was created. | 
**UpdatedAt** | **time.Time** | The timestamp when the employee record was last updated. | 

## Methods

### NewEmployeeAttributes

`func NewEmployeeAttributes(distributorId uuid.UUID, role string, createdAt time.Time, updatedAt time.Time, ) *EmployeeAttributes`

NewEmployeeAttributes instantiates a new EmployeeAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmployeeAttributesWithDefaults

`func NewEmployeeAttributesWithDefaults() *EmployeeAttributes`

NewEmployeeAttributesWithDefaults instantiates a new EmployeeAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDistributorId

`func (o *EmployeeAttributes) GetDistributorId() uuid.UUID`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *EmployeeAttributes) GetDistributorIdOk() (*uuid.UUID, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *EmployeeAttributes) SetDistributorId(v uuid.UUID)`

SetDistributorId sets DistributorId field to given value.


### GetRole

`func (o *EmployeeAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *EmployeeAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *EmployeeAttributes) SetRole(v string)`

SetRole sets Role field to given value.


### GetCreatedAt

`func (o *EmployeeAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *EmployeeAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *EmployeeAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *EmployeeAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *EmployeeAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *EmployeeAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


