# CreateEmployeeInviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CompanyId** | [**uuid.UUID**](uuid.UUID.md) | The UUID of the company to which the user is being invited. | 
**UserId** | [**uuid.UUID**](uuid.UUID.md) | The UUID of the user being invited. | 
**Role** | **string** | The role assigned to the invited user within the company. | 

## Methods

### NewCreateEmployeeInviteDataAttributes

`func NewCreateEmployeeInviteDataAttributes(companyId uuid.UUID, userId uuid.UUID, role string, ) *CreateEmployeeInviteDataAttributes`

NewCreateEmployeeInviteDataAttributes instantiates a new CreateEmployeeInviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateEmployeeInviteDataAttributesWithDefaults

`func NewCreateEmployeeInviteDataAttributesWithDefaults() *CreateEmployeeInviteDataAttributes`

NewCreateEmployeeInviteDataAttributesWithDefaults instantiates a new CreateEmployeeInviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompanyId

`func (o *CreateEmployeeInviteDataAttributes) GetCompanyId() uuid.UUID`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *CreateEmployeeInviteDataAttributes) GetCompanyIdOk() (*uuid.UUID, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *CreateEmployeeInviteDataAttributes) SetCompanyId(v uuid.UUID)`

SetCompanyId sets CompanyId field to given value.


### GetUserId

`func (o *CreateEmployeeInviteDataAttributes) GetUserId() uuid.UUID`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *CreateEmployeeInviteDataAttributes) GetUserIdOk() (*uuid.UUID, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *CreateEmployeeInviteDataAttributes) SetUserId(v uuid.UUID)`

SetUserId sets UserId field to given value.


### GetRole

`func (o *CreateEmployeeInviteDataAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *CreateEmployeeInviteDataAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *CreateEmployeeInviteDataAttributes) SetRole(v string)`

SetRole sets Role field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


