# EmployeeInviteAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**companyId** | **string** | The unique identifier for the company associated with the employee invite. | 
**UserId** | **string** | The unique identifier for the user being invited as an employee. | 
**InvitedBy** | **string** | The unique identifier for the user who sent the invitation. | 
**Role** | **string** | The role of the invited employee within the company&#39;s organization (e.g., manager, staff). | 
**Status** | **string** | The current status of the employee invitation. | 
**AnsweredAt** | Pointer to **time.Time** | The timestamp when the invitation was responded to (accepted or declined). | [optional] 
**CreatedAt** | **time.Time** | The timestamp when the employee invite was created. | 

## Methods

### NewEmployeeInviteAttributes

`func NewEmployeeInviteAttributes(companyId string, userId string, invitedBy string, role string, status string, createdAt time.Time, ) *EmployeeInviteAttributes`

NewEmployeeInviteAttributes instantiates a new EmployeeInviteAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmployeeInviteAttributesWithDefaults

`func NewEmployeeInviteAttributesWithDefaults() *EmployeeInviteAttributes`

NewEmployeeInviteAttributesWithDefaults instantiates a new EmployeeInviteAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetcompanyId

`func (o *EmployeeInviteAttributes) GetcompanyId() string`

GetcompanyId returns the companyId field if non-nil, zero value otherwise.

### GetcompanyIdOk

`func (o *EmployeeInviteAttributes) GetcompanyIdOk() (*string, bool)`

GetcompanyIdOk returns a tuple with the companyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetcompanyId

`func (o *EmployeeInviteAttributes) SetcompanyId(v string)`

SetcompanyId sets companyId field to given value.


### GetUserId

`func (o *EmployeeInviteAttributes) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *EmployeeInviteAttributes) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *EmployeeInviteAttributes) SetUserId(v string)`

SetUserId sets UserId field to given value.


### GetInvitedBy

`func (o *EmployeeInviteAttributes) GetInvitedBy() string`

GetInvitedBy returns the InvitedBy field if non-nil, zero value otherwise.

### GetInvitedByOk

`func (o *EmployeeInviteAttributes) GetInvitedByOk() (*string, bool)`

GetInvitedByOk returns a tuple with the InvitedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvitedBy

`func (o *EmployeeInviteAttributes) SetInvitedBy(v string)`

SetInvitedBy sets InvitedBy field to given value.


### GetRole

`func (o *EmployeeInviteAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *EmployeeInviteAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *EmployeeInviteAttributes) SetRole(v string)`

SetRole sets Role field to given value.


### GetStatus

`func (o *EmployeeInviteAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *EmployeeInviteAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *EmployeeInviteAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetAnsweredAt

`func (o *EmployeeInviteAttributes) GetAnsweredAt() time.Time`

GetAnsweredAt returns the AnsweredAt field if non-nil, zero value otherwise.

### GetAnsweredAtOk

`func (o *EmployeeInviteAttributes) GetAnsweredAtOk() (*time.Time, bool)`

GetAnsweredAtOk returns a tuple with the AnsweredAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnsweredAt

`func (o *EmployeeInviteAttributes) SetAnsweredAt(v time.Time)`

SetAnsweredAt sets AnsweredAt field to given value.

### HasAnsweredAt

`func (o *EmployeeInviteAttributes) HasAnsweredAt() bool`

HasAnsweredAt returns a boolean if a field has been set.

### GetCreatedAt

`func (o *EmployeeInviteAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *EmployeeInviteAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *EmployeeInviteAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


