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