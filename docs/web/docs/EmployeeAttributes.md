# EmployeeAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Role** | **string** | The role of the employee within the company&#39;s organization (e.g., manager, staff). | 
**Position** | Pointer to **string** | The job position or title of the employee. | [optional] 
**Label** | Pointer to **string** | A human-readable label or name for the employee. | [optional] 
**CreatedAt** | **time.Time** | The timestamp when the employee record was created. | 
**UpdatedAt** | **time.Time** | The timestamp when the employee record was last updated. | 

## Methods

### NewEmployeeAttributes

`func NewEmployeeAttributes(role string, createdAt time.Time, updatedAt time.Time, ) *EmployeeAttributes`

NewEmployeeAttributes instantiates a new EmployeeAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmployeeAttributesWithDefaults

`func NewEmployeeAttributesWithDefaults() *EmployeeAttributes`

NewEmployeeAttributesWithDefaults instantiates a new EmployeeAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

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


### GetPosition

`func (o *EmployeeAttributes) GetPosition() string`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *EmployeeAttributes) GetPositionOk() (*string, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *EmployeeAttributes) SetPosition(v string)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *EmployeeAttributes) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### GetLabel

`func (o *EmployeeAttributes) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *EmployeeAttributes) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *EmployeeAttributes) SetLabel(v string)`

SetLabel sets Label field to given value.

### HasLabel

`func (o *EmployeeAttributes) HasLabel() bool`

HasLabel returns a boolean if a field has been set.

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


