# sunny-boy-bigquery-exporter
This application exports readings from an SMA sunny boy PV inverter into BigQuery

## ModBus registers for SM SB 4000TL-21

| Modbus register number | SunSpec-Name | Description / Number code(s) | Type | Access |
| ---------------------- | ------------ | ---------------------------- | ---- | ------ |
| 40001 | SID | A well-known value 0x53756e53.  Uniquely identifies this as a SunSpec Modbus Map: 1400204883 | uint32 | RO |
| 40003 | ID | A well-known value 1.  Uniquely identifies this as a SunSpec Common Model | uint16 | RO |
| 40004 | L | Well-known # of 16 bit registers to follow : 66 | uint16 | RO |
| 40005 | Mn | "Well known value registered with SunSpec for compliance:
SMA" | string16 | RO |
| 40021 | Md | "Manufacturer specific value (32 chars):
Solar Inverter" | string16 | RO |
| 40037 | Opt | Manufacturer specific value (16 chars): Model ID | string8 | RO |
| 40045 | Vr | Manufacturer specific value (16 chars) | string8 | RO |
| 40053 | SN | Manufacturer specific value (32 chars) | string16 | RO |
| 40071 | ID | A well-known value 11.  Uniquely identifies this as a SunSpec Ethernet Link Layer Model | uint16 | RO |
| 40072 | L | Well-known # of 16 bit registers to follow : 13 | uint16 | RO |
| 40073 | Spd | "Interface speed in Mb/s:
0
10
100" | uint16 | RO |
| 40074 | CfgSt | "Bitmask values Interface flags:
0
1
3" | uint16 | RO |
| 40075 | St | "Enumerated value. State information for this interface:
0
1
2" | uint16 | RO |
| 40076 | MAC | IEEE MAC address of this interface | uint64 | RO | 
| 40086 | ID | A well-known value 12.  Uniquely identifies this as a SunSpec IPv4 Model | uint16 | RO |
| 40087 | L | Well-known # of 16 bit registers to follow : 98 | uint16 | RO |
| 40092 | CfgSt | Enumerated value.  Configuration status | uint16 | RO |
| 40093 | ChgSt | Bitmask value.  A configuration change is pending | uint16 | RO |
| 40094 | Cap | Bitmask value. Identify capable sources of configuration | uint16 | RO |
| 40095 | Cfg | "Enumerated value.  Configuration method used:
0
1" | uint16 | RW
| 40096 | Ctl | Configure use of services | uint16 | RO |
| 40097 | Addr | IPv4 numeric address as a dotted string xxx.xxx.xxx.xxx | string8 | RW
| 40105 | Msk | IPv4 numeric netmask as a dotted string xxx.xxx.xxx.xxx | string8 | RW
| 40113 | Gw | IPv4 numeric gateway address as a dotted string xxx.xxx.xxx.xxx | string8 | RW
| 40121 | DNS1 | IPv4 numeric DNS address as a dotted string xxx.xxx.xxx.xxx | string8 | RW | 
| 40186 | ID | "A well-known value 101,102,103.  Uniquely identifies this as a SunSpec Inverter Model (101: 1-phase, 102: 2-phase, 103: 3-phase):
101
103" | uint16 | RO |
| 40187 | L | Well-known # of 16 bit registers to follow : 50 | uint16 | RO |
| 40189 | AphA | Phase A Current | uint16 | RO |
| 40190 | AphB | Phase B Current | uint16 | RO |
| 40191 | AphC | Phase C Current | uint16 | RO |
| 40192 | A_SF | int16 | RO |
| 40196 | PhVphA | Phase Voltage AN | uint16 | RO |
| 40197 | PhVphB | Phase Voltage BN | uint16 | RO |
| 40198 | PhVphC | Phase Voltage CN | uint16 | RO |
| 40199 | V_SF | int16 | RO |
| 40200 | W | AC Power | int16 | RO |
| 40201 | W_SF | int16 | RO |
| 40202 | Hz | Line Frequency | uint16 | RO |
| 40203 | Hz_SF | int16 | RO |
| 40204 | VA | AC Apparent Power | int16 | RO |
| 40205 | VA_SF | int16 | RO |
| 40206 | VAr | AC Reactive Power | int16 | RO |
| 40207 | VAr_SF | int16 | RO |
| 40208 | PF | AC Power Factor | int16 | RO |
| 40209 | PF_SF | int16 | RO |
| 40210 | WH | AC Energy | acc32 | RO |
| 40212 | WH_SF | int16 | RO |
| 40218 | DCW_SF | int16 | RO |
| 40223 | Tmp_SF | int16 | RO |
| 40226 | Evt1 | "Bitmask value. Event fields:
0
1
2
16
32
64
128
256
512
1024
2048
4096
8192
16384
32768" | uint32 | RO |
| 40238 | ID | A well-known value 120.  Uniquely identifies this as a SunSpec Nameplate Model | uint16 | RO |
| 40239 | L | Well-known # of 16 bit registers to follow : 26 | uint16 | RO |
| 40240 | DERTyp | "Type of DER device. Default value is 4 to indicate PV device:
4" | uint16 | RO |
| 40241 | WRtg | Continuous power output capability of the inverter | uint16 | RO |
| 40242 | WRtg_SF | Scale factor | int16 | RO |
| 40243 | VARtg | Continuous Volt-Ampere capability of the inverter | uint16 | RO |
| 40244 | VARtg_SF | Scale factor | int16 | RO |
| 40245 | VArRtgQ1 | Continuous VAR capability of the inverter in quadrant 1 | int16 | RO |
| 40246 | VArRtgQ2 | Continuous VAR capability of the inverter in quadrant 2 | int16 | RO |
| 40247 | VArRtgQ3 | Continuous VAR capability of the inverter in quadrant 3 | int16 | RO |
| 40248 | VArRtgQ4 | Continuous VAR capability of the inverter in quadrant 4 | int16 | RO |
| 40249 | VArRtg_SF | Scale factor | int16 | RO |
| 40251 | ARtg_SF | Scale factor | int16 | RO |
| 40252 | PFRtgQ1 | Minimum power factor capability of the inverter in quadrant 1 | int16 | RO |
| 40253 | PFRtgQ2 | Minimum power factor capability of the inverter in quadrant 2 | int16 | RO |
| 40254 | PFRtgQ3 | Minimum power factor capability of the inverter in quadrant 3 | int16 | RO |
| 40255 | PFRtgQ4 | Minimum power factor capability of the inverter in quadrant 4 | int16 | RO |
| 40256 | PFRtg_SF | Scale factor | int16 | RO |
| 40258 | WHRtg_SF | Scale factor | int16 | RO |
| 40260 | AhrRtg_SF | Scale factor for amp-hour rating | int16 | RO |
| 40262 | MaxChaRte_SF | Scale factor | int16 | RO |
| 40264 | MaxDisChaRte_SF | Scale factor | int16 | RO |
|  | 
| 40266 | ID | A well-known value 121.  Uniquely identifies this as a SunSpec Basic Settings Model | uint16 | RO |
| 40267 | L | Well-known # of 16 bit registers to follow : 30 | uint16 | RO |
| 40268 | WMax | Setting for maximum power output. Default to WRtg | uint16 | RW
| 40269 | VRef | Voltage at the PCC | uint16 | RW
| 40270 | VRefOfs | Offset  from PCC to inverter | int16 | RW
| 40273 | VAMax | Setpoint for maximum apparent power. Default to VARtg | uint16 | RW
| 40278 | WGra | Default ramp rate of change of active power due to command or internal action | uint16 | RW
| 40286 | ECPNomHz | Setpoint for nominal frequency at the ECP | uint16 | RW
| 40287 | ConnPh | "Identity of connected phase for single phase inverters. A=1 B=2 C=3:
0
1
2
3" | uint16 | RW
| 40288 | WMax_SF | Scale factor for real power | int16 | RO |
| 40289 | VRef_SF | Scale factor for voltage at the PCC | int16 | RO |
| 40290 | VRefOfs_SF | Scale factor for offset voltage | int16 | RO |
| 40291 | VMinMax_SF | Scale factor for min/max voltages | int16 | RO |
| 40292 | VAMax_SF | Scale factor for apparent power | int16 | RO |
| 40293 | VArMax_SF | Scale factor for reactive power | int16 | RO |
| 40294 | WGra_SF | Scale factor for default ramp rate | int16 | RO |
| 40297 | ECPNomHz_SF | Scale factor for nominal frequency | int16 | RO |
 | 
| 40298 | ID | A well-known value 122.  Uniquely identifies this as a SunSpec Measurements_Status Model | uint16 | RO |
| 40299 | L | Well-known # of 16 bit registers to follow : 44 | uint16 | RO |
| 40300 | PVConn | "PV inverter present/available status. Enumerated value:
1
3
5" | uint16 | RO |
| 40301 | StorConn | "Storage inverter present/available status. Enumerated value:
0" | uint16 | RO |
| 40302 | ECPConn | "ECP connection status: disconnected=0  connected=1:
0" | uint16 | RO |
| 40303 | ActWh | AC lifetime active (real) energy output | acc64 | RO |
| 40342 | Ris | Isolation resistance | uint16 | RO |
| 40343 | Ris_SF | Scale factor for isolation resistance | int16 | RO |
 | 
| 40344 | ID | A well-known value 123.  Uniquely identifies this as a SunSpec Immediate Controls Model | uint16 | RO |
| 40345 | L | Well-known # of 16 bit registers to follow : 24 | uint16 | RO |
| 40353 | WMaxLim_Ena | "Enumerated valued.  Throttle enable/disable control:
0
1" | uint16 | RW
| 40358 | OutPFSet_Ena | "Enumerated valued.  Fixed power factor enable/disable control:
0
1" | uint16 | RW
| 40365 | VArPct_Mod | Enumerated value. VAR percent limit mode | uint16 | RO |
| 40366 | VArPct_Ena | "Enumerated valued.  Percent limit VAr enable/disable control:
0
1" | uint16 | RW
| 40367 | WMaxLimPct_SF | Scale factor for power output percent | int16 | RO |
| 40368 | OutPFSet_SF | Scale factor for power factor | int16 | RO |
| 40369 | VArPct_SF | Scale factor for reactive power percent | int16 | RO |
 | 
| 40370 | ID | A well-known value 124.  Uniquely identifies this as a SunSpec Storage Model | uint16 | RO |
| 40371 | L | Well-known # of 16 bit registers to follow : 24 | uint16 | RO |
| 40375 | StorCtl_Mod | "Activate hold/discharge/charge storage control mode. Bitfield value:
0" | uint16 | RW
| 40388 | WChaMax_SF | Scale factor for maximum charge | int16 | RO |
| 40392 | ChaState_SF | Scale factor for available energy percent | int16 | RO |
| 40394 | InBatV_SF | Scale factor for battery voltage | int16 | RO |
 | 
| 40396 | ID | A well-known value 126.  Uniquely identifies this as a SunSpec Static Volt-VAR Model | uint16 | RO |
| 40397 | L | Variable # of 16 bit registers to follow : 10+N*54 | uint16 | RO |
| 40398 | ActCrv | "Index of active curve. 0=no active curve:
0" | uint16 | RW
| 40399 | ModEna | "Is Volt-VAR control active:
0" | uint16 | RW
| 40403 | NCrv | Number of curves supported (recommend 4) | uint16 | RO |
| 40404 | NPt | Number of curve points supported (maximum of 20) | uint16 | RO |
| 40405 | V_SF | Scale factor for percent VRef | int16 | RO |
| 40406 | DeptRef_SF | scale factor for dependent variable | int16 | RO |
| 40407 | RmpIncDec_SF | Scale factor for increment and decrement ramps | int16 | RO |
| 40409 | DeptRef | Meaning of dependent variable: 1=%WMax 2=%VArMax 3=%VArAval | uint16 | RO |
| 40461 | ReadOnly | Boolean flag indicates if curve is read-only or can be modified | uint16 | RO |
 | 
| 40462 | ID | A well-known value 127.  Uniquely identifies this as a SunSpec Freq-Watt Param Model | uint16 | RO |
| 40463 | L | Well-known # of 16 bit registers to follow : 10 | uint16 | RO |
| 40464 | WGra | The slope of the reduction in the maximum allowed watts output as a function of frequency | uint16 | RW
| 40465 | HzStr | The frequency deviation from nominal frequency (ECPNomHz) at which a snapshot of the instantaneous power output is taken to act as the CAPPED power level (PM) and above which reduction in power output occurs | int16 | RW
| 40466 | HzStop | The frequency deviation from nominal frequency (ECPNomHz) at which curtailed power output may return to normal and the cap on the power level value is removed | int16 | RW
| 40467 | HysEna | "Enable hysterisis:
0
1" | uint16 | RW
| 40468 | ModEna | "Is Parameterized Frequency-Watt control active:
0
1" | uint16 | RW
| 40469 | HzStopWGra | The maximum time-based rate of change at which power output returns to normal after having been capped by an over frequency event | uint16 | RW
| 40470 | WGra_SF | Scale factor for output gradient | int16 | RO |
| 40471 | HzStrStop_SF | Scale factor for frequency deviations | int16 | RO |
| 40472 | RmpIncDec_SF | Scale factor for increment and decrement ramps | int16 | RO |
 | 
| 40474 | ID | A well-known value 128.  Uniquely identifies this as a SunSpec Dynamic Reactive Current Model | uint16 | RO |
| 40475 | L | Well-known # of 16 bit registers to follow : 14 | uint16 | RO |
| 40479 | ModEna | "Activate dynamic reactive current model:
0" | uint16 | RW
| 40480 | FilTms | The time window used to calculate the moving average voltage | uint16 | RO |
| 40483 | BlkZnV | Block zone voltage which defines a lower voltage boundary below which no dynamic current is produced | uint16 | RW
| 40487 | ArGra_SF | Scale factor for the gradients | int16 | RO |
| 40488 | VRefPct_SF | Scale factor for the voltage zone and limit settings | int16 | RO |
 | 
| 40490 | ID | A well-known value 131.  Uniquely identifies this as a SunSpec Watt-PF Model | uint16 | RO |
| 40491 | L | Variable # of 16 bit registers to follow : 10+N*54 | uint16 | RO |
| 40492 | ActCrv | "Index of active curve. 0=no active curve:
0
1" | uint16 | RW
| 40493 | ModEna | "Is watt-PF mode active:
0
1" | uint16 | RW
| 40498 | NPt | Max number of points in array | uint16 | RO |
| 40499 | W_SF | Scale factor for percent WMax | int16 | RO |
| 40500 | PF_SF | Scale factor for PF | int16 | RO |
| 40502 | ActPt | Number of active points in array | uint16 | RO |
| 40503 | W1 | Point 1 Watts | int16 | RW
| 40505 | W2 | Point 2 Watts | int16 | RW
| 40554 | ReadOnly | Enumerated value indicates if curve is read-only or can be modified | uint16 | RO |
 | 
| 40556 | ID | A well-known value 132.  Uniquely identifies this as a SunSpec Volt-Watt Model | uint16 | RO |
| 40557 | L | Variable # of 16 bit registers to follow : 10+N*54 | uint16 | RO |
| 40558 | ActCrv | "Index of active curve. 0=no active curve:
0
1" | uint16 | RW
| 40559 | ModEna | "Is Volt-Watt control active:
0
1" | uint16 | RW
| 40563 | NCrv | Number of curves supported (recommend min. 4) | uint16 | RO |
| 40564 | NPt | Number of points in array (maximum 20) | uint16 | RO |
| 40565 | V_SF | Scale factor for percent VRef | int16 | RO |
| 40566 | DeptRef_SF | Scale Factor for % DeptRef | int16 | RO |
| 40569 | DeptRef | Defines the meaning of the Watts DeptRef.  1=% WMax 2=% WAvail | uint16 | RO |
| 40621 | ReadOnly | Enumerated value indicates if curve is read-only or can be modified | uint16 | RO |
 | 
| 40622 | ID | A well-known value 160.  Uniquely identifies this as a SunSpec Multiple MPPT Inverter Extension Model Model | uint16 | RO |
| 40623 | L | Variable # of 16 bit registers to follow : 8+N*20 | uint16 | RO |
| 40624 | DCA_SF | Current Scale Factor | int16 | RO |
| 40625 | DCV_SF | Voltage Scale Factor | int16 | RO |
| 40626 | DCW_SF | Power Scale Factor | int16 | RO |
| 40630 | N | Number of Modules | uint16 | RO |
| 40632 | ID | MPPT 1: Input ID | uint16 | RO |
| 40641 | DCA | MPPT 1: DC Current | uint16 | RO |
| 40642 | DCV | MPPT 1: DC Voltage | uint16 | RO |
| 40643 | DCW | MPPT 1: DC Power | uint16 | RO |
| 40652 | ID | MPPT 2: Input ID | uint16 | RO |
| 40661 | DCA | MPPT 2: DC Current | uint16 | RO |
| 40662 | DCV | MPPT 2: DC Voltage | uint16 | RO |
| 40663 | DCW | MPPT 2: DC Power | uint16 | RO |