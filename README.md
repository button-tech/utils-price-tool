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
    "change": "24", 
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
                    "0x0000000000000000000000000000000000000000": "8599.95286058",
                    "percent_change": "-12.52"
                },
                {
                    "0x000000000000000000000000000000000000003C": "166.238229345",
                    "percent_change": "-19.11"
                }
            ]
        },
        {
            "currency": "RUB",
            "rates": [
                {
                    "0x0000000000000000000000000000000000000000": "547576.1985388497",
                    "percent_change": "-12.52"
                },
                {
                    "0x000000000000000000000000000000000000003C": "10584.72053885484",
                    "percent_change": "-19.11"
                }
            ]
        }
    ]
}
```

API can be:

```"cmc"``` - Coin-market-cap, 

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