# CreateCompanyBlockDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CompanyId** | [**uuid.UUID**](uuid.UUID.md) | The UUID of the company to be blocked. | 
**Reason** | **string** | The reason for blocking the company. | 

## Methods

### NewCreateCompanyBlockDataAttributes

`func NewCreateCompanyBlockDataAttributes(companyId uuid.UUID, reason string, ) *CreateCompanyBlockDataAttributes`

NewCreateCompanyBlockDataAttributes instantiates a new CreateCompanyBlockDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateCompanyBlockDataAttributesWithDefaults

`func NewCreateCompanyBlockDataAttributesWithDefaults() *CreateCompanyBlockDataAttributes`

NewCreateCompanyBlockDataAttributesWithDefaults instantiates a new CreateCompanyBlockDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompanyId

`func (o *CreateCompanyBlockDataAttributes) GetCompanyId() uuid.UUID`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *CreateCompanyBlockDataAttributes) GetCompanyIdOk() (*uuid.UUID, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *CreateCompanyBlockDataAttributes) SetCompanyId(v uuid.UUID)`

SetCompanyId sets CompanyId field to given value.


### GetReason

`func (o *CreateCompanyBlockDataAttributes) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *CreateCompanyBlockDataAttributes) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *CreateCompanyBlockDataAttributes) SetReason(v string)`

SetReason sets Reason field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


