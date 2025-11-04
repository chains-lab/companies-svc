# CreateInviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserId** | [**uuid.UUID**](uuid.UUID.md) | ID of the user being invited | 
**CompanyId** | [**uuid.UUID**](uuid.UUID.md) | ID of the company the invite is for | 
**Role** | **string** | Role assigned to the invited user | 

## Methods

### NewCreateInviteDataAttributes

`func NewCreateInviteDataAttributes(userId uuid.UUID, companyId uuid.UUID, role string, ) *CreateInviteDataAttributes`

NewCreateInviteDataAttributes instantiates a new CreateInviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateInviteDataAttributesWithDefaults

`func NewCreateInviteDataAttributesWithDefaults() *CreateInviteDataAttributes`

NewCreateInviteDataAttributesWithDefaults instantiates a new CreateInviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserId

`func (o *CreateInviteDataAttributes) GetUserId() uuid.UUID`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *CreateInviteDataAttributes) GetUserIdOk() (*uuid.UUID, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *CreateInviteDataAttributes) SetUserId(v uuid.UUID)`

SetUserId sets UserId field to given value.


### GetCompanyId

`func (o *CreateInviteDataAttributes) GetCompanyId() uuid.UUID`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *CreateInviteDataAttributes) GetCompanyIdOk() (*uuid.UUID, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *CreateInviteDataAttributes) SetCompanyId(v uuid.UUID)`

SetCompanyId sets CompanyId field to given value.


### GetRole

`func (o *CreateInviteDataAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *CreateInviteDataAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *CreateInviteDataAttributes) SetRole(v string)`

SetRole sets Role field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


