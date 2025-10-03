# CreateDistributorBlockDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DistributorId** | [**uuid.UUID**](uuid.UUID.md) | The UUID of the distributor to be blocked. | 
**Reason** | **string** | The reason for blocking the distributor. | 

## Methods

### NewCreateDistributorBlockDataAttributes

`func NewCreateDistributorBlockDataAttributes(distributorId uuid.UUID, reason string, ) *CreateDistributorBlockDataAttributes`

NewCreateDistributorBlockDataAttributes instantiates a new CreateDistributorBlockDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateDistributorBlockDataAttributesWithDefaults

`func NewCreateDistributorBlockDataAttributesWithDefaults() *CreateDistributorBlockDataAttributes`

NewCreateDistributorBlockDataAttributesWithDefaults instantiates a new CreateDistributorBlockDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDistributorId

`func (o *CreateDistributorBlockDataAttributes) GetDistributorId() uuid.UUID`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *CreateDistributorBlockDataAttributes) GetDistributorIdOk() (*uuid.UUID, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *CreateDistributorBlockDataAttributes) SetDistributorId(v uuid.UUID)`

SetDistributorId sets DistributorId field to given value.


### GetReason

`func (o *CreateDistributorBlockDataAttributes) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *CreateDistributorBlockDataAttributes) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *CreateDistributorBlockDataAttributes) SetReason(v string)`

SetReason sets Reason field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


