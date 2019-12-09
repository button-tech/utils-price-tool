This API allows you to get current exchange rates from crypto-compare and trust wallet.
 
The main feature is support of trust-wallet requests.

```
                                    Version v1
```
Method POST:

```/courses/v2/prices```
Example:

```Request body:```

```
{
    "tokens": [
        "0", "457"
    ],
    "currencies": [
        "USD"
    ],
    "change": "1",
    "api": "pcmc"
}
```

```Response```

```
{
    "data": [
        {
            "currency": "USD",
            "rates": [
                {
                    "0": "7469.7828042500",
                    "percent_change": "-0.005513"
                },
                {
                    "457": "0.1638779609",
                    "percent_change": "0.539258"
                }
            ]
        }
    ]
}
```

Method GET:

```/courses/v2/info```

```
{
    "api": [
        {
            "name": "ntrust",
            "supported_changes": [
                "0",
                "24"
            ],
            "supported_fiats": {
                "AE": 457,
                "AION": 425,
                "ALGO": 283,
                "ATOM": 118,
                "BCH": 145,
                "BNB": 714,
                "BTC": 0,
                "CLO": 820,
                "DASH": 5,
                "DCR": 42,
                "DGB": 20,
                "DOGE": 3,
                "ETC": 61,
                "ETH": 60,
                "FIO": 235,
                "GO": 6060,
                "GRS": 17,
                "ICX": 74,
                "IOTX": 304,
                "KIN": 2017,
                "LTC": 2,
                "NANO": 165,
                "NAS": 2718,
                "NIM": 242,
                "ONT": 1024,
                "POA": 178,
                "QTUM": 2301,
                "RVN": 175,
                "SEM": 7562605,
                "THETA": 500,
                "TOMO": 889,
                "TRX": 195,
                "TT": 1001,
                "VET": 818,
                "VIA": 14,
                "WAN": 5718350,
                "WAVES": 5741564,
                "XLM": 148,
                "XRP": 144,
                "XTZ": 1729,
                "XZC": 136,
                "ZEC": 133,
                "ZEL": 19167,
                "ZIL": 313
            }
        },
        {
            "name": "pcmc",
            "supported_changes": [
                "0",
                "24",
                "7d"
            ],
            "supported_fiats": {
                "AE": 457,
                "ALGO": 283,
                "ATOM": 118,
                "BCH": 145,
                "BNB": 714,
                "BTC": 0,
                "DASH": 5,
                "DCR": 42,
                "DGB": 20,
                "DOGE": 3,
                "ETC": 61,
                "ETH": 60,
                "ICX": 74,
                "LTC": 2,
                "NANO": 165,
                "ONT": 1024,
                "QTUM": 2301,
                "RVN": 175,
                "THETA": 500,
                "TRX": 195,
                "VET": 818,
                "WAVES": 5741564,
                "XLM": 148,
                "XRP": 144,
                "XTZ": 1729,
                "ZEC": 133,
                "ZIL": 313
            }
        }
    ]
}
```

```
                                    Version v1
```

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

```"cmc"``` - Coin-market-cap

```"crc""``` - Crypto-compare

```"huobi""``` - Huobi.pro


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
        },
        {
            "name": "huobi",
            "supported_changes": [
                "0"
            ]
        }
    ]
}
```


API supports this list of fiats:
```
        "AED",
	"AFN",
	"ALL",
	"AMD",
	"ANG",
	"AOA",
	"ARS",
	"AUD",
	"AWG",
	"AZN",
	"BAM",
	"BBD",
	"BDT",
	"BGN",
	"BHD",
	"BIF",
	"BMD",
	"BND",
	"BOB",
	"BRL",
	"BSD",
	"BTC",
	"BTN",
	"BWP",
	"BYN",
	"BYR",
	"BZD",
	"CAD",
	"CDF",
	"CHF",
	"CLF",
	"CLP",
	"CNY",
	"COP",
	"CRC",
	"CUC",
	"CUP",
	"CVE",
	"CZK",
	"DJF",
	"DKK",
	"DOP",
	"DZD",
	"EGP",
	"ERN",
	"ETB",
	"EUR",
	"FJD",
	"FKP",
	"GBP",
	"GEL",
	"GGP",
	"GHS",
	"GIP",
	"GMD",
	"GNF",
	"GTQ",
	"GYD",
	"HKD",
	"HNL",
	"HRK",
	"HTG",
	"HUF",
	"IDR",
	"ILS",
	"IMP",
	"INR",
	"IQD",
	"IRR",
	"ISK",
	"JEP",
	"JMD",
	"JOD",
	"JPY",
	"KES",
	"KGS",
	"KHR",
	"KMF",
	"KPW",
	"KRW",
	"KWD",
	"KYD",
	"KZT",
	"LAK",
	"LBP",
	"LKR",
	"LRD",
	"LSL",
	"LTL",
	"LVL",
	"LYD",
	"MAD",
	"MDL",
	"MGA",
	"MKD",
	"MMK",
	"MNT",
	"MOP",
	"MRO",
	"MUR",
	"MVR",
	"MWK",
	"MXN",
	"MYR",
	"MZN",
	"NAD",
	"NGN",
	"NIO",
	"NOK",
	"NPR",
	"NZD",
	"OMR",
	"PAB",
	"PEN",
	"PGK",
	"PHP",
	"PKR",
	"PLN",
	"PYG",
	"QAR",
	"RON",
	"RUB",
	"RWF",
	"SAR",
	"SBD",
	"SCR",
	"SDG",
	"SEK",
	"SGD",
	"SHP",
	"SLL",
	"SOS",
	"SRD",
	"STD",
	"SVC",
	"SYP",
	"SZL",
	"THB",
	"TJS",
	"TMT",
	"TND",
	"TOP",
	"TRY",
	"TTD",
	"TWD",
	"TZS",
	"UAH",
	"UGX",
	"USD",
	"UYU",
	"UZS",
	"VEF",
	"VND",
	"VUV",
	"WST",
	"XAF",
	"XAG",
	"XAU",
	"XCD",
	"XDR",
	"XOF",
	"XPF",
	"YER",
	"ZAR",
	"ZMK",
	"ZMW",
	"ZWL",
```

But, supported API don't work with all fiat.

Soon there will be a list of specific supported API with their supported fiats.