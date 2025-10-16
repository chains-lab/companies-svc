# EmployeeAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | **string** | The username of the employee within the company&#39;s system. | 
**Avatar** | **string** | A URL pointing to the employee&#39;s avatar image. | 
**CompanyId** | [**uuid.UUID**](uuid.UUID.md) | The unique identifier for the company associated with the employee. | 
**Role** | **string** | The role of the employee within the company&#39;s organization (e.g., manager, staff). | 
**CreatedAt** | **time.Time** | The timestamp when the employee record was created. | 
**UpdatedAt** | **time.Time** | The timestamp when the employee record was last updated. | 

## Methods

### NewEmployeeAttributes

`func NewEmployeeAttributes(username string, avatar string, companyId uuid.UUID, role string, createdAt time.Time, updatedAt time.Time, ) *EmployeeAttributes`

NewEmployeeAttributes instantiates a new EmployeeAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmployeeAttributesWithDefaults

`func NewEmployeeAttributesWithDefaults() *EmployeeAttributes`

NewEmployeeAttributesWithDefaults instantiates a new EmployeeAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *EmployeeAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *EmployeeAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *EmployeeAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetAvatar

`func (o *EmployeeAttributes) GetAvatar() string`

GetAvatar returns the Avatar field if non-nil, zero value otherwise.

### GetAvatarOk

`func (o *EmployeeAttributes) GetAvatarOk() (*string, bool)`

GetAvatarOk returns a tuple with the Avatar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatar

`func (o *EmployeeAttributes) SetAvatar(v string)`

SetAvatar sets Avatar field to given value.


### GetCompanyId

`func (o *EmployeeAttributes) GetCompanyId() uuid.UUID`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *EmployeeAttributes) GetCompanyIdOk() (*uuid.UUID, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *EmployeeAttributes) SetCompanyId(v uuid.UUID)`

SetCompanyId sets CompanyId field to given value.


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


