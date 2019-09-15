This API allows you to get current exchange rates from crypto-compare and trust wallet.
 
The main feature is support of trust-wallet requests.

Method POST: 

```/courses/v1/prices```

Examples:
```
Request body:
```

```
{
    "tokens": [
        "0x0000000000000000000000000000000000000000",
        "0x000000000000000000000000000000000000003C"
    ],
    "currencies": [
        "USD",
        "RUB"
    ],
    "change": "0", - you can skip, because value is 0
    "api": "cmc"
}
```
```
Response:
```
```
{
    "data": [
        {
            "currency": "USD",
            "rates": [
                {
                    "rate": {
                        "0x0000000000000000000000000000000000000000": "10338.6989468"
                    },
                    "percent_change": "-0.04"
                },
                {
                    "rate": {
                        "0x000000000000000000000000000000000000003C": "188.300781866"
                    },
                    "percent_change": "3.98"
                }
            ]
        }
    ]
}
```

API can be:

```"cmc"``` - Coinmarketcap, 

```"crc""``` - Crypto-compare.


Method GET send you data of all confirmed API and changes: 

```/courses/v1/list```

```
Response:

{
    "api": [
        {
            "name": "crc",
            "supported_changes": [
                "0",
                "1",
                "24"
            ]
        },
        {
            "name": "cmc",
            "supported_changes": [
                "24"
            ]
        }
    ]
}
```