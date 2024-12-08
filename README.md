# Landlord API
Landlord API is useful for calculation of a field when land lord need to build an apartment there. 

It roughly calculates the building area and how many lots/flats can be yield from that land. 


## Endpoints
- POST /land/calculate endpoint is taking input as properties of land and allowed rates for making an apartment there. In the response it returns how many 1+1 or 2+1 or 3+1 flat can be built on that land.

- Request:
```json
{
    "emsal": 0.8,
    "landSquareMeter": 600,
    "maxAllowedFlatCount": 5,
    "extendFactor": 1.3
}
```
- Response:
```json
{
    "OnePlusOneCount": [
        {
            "Count": 17,
            "Size": 36.705882352941174
        },
        {
            "Count": 14,
            "Size": 44.57142857142857
        },
        {
            "Count": 13,
            "Size": 48
        }
    ],
    "TwoPlusOneCount": 0,
    "ThreePlusOneCount": 0
}
```