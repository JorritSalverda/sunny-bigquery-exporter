package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// SunnyBoyClient is the interface for connecting to sma sunny boy inverter via ethernet
type SunnyBoyClient interface {
	ReadInputRegisters(address, quantity uint16) (results []byte, err error)
	Login() (success bool, err error)
}

type sunnyBoyClientImpl struct {
	localPort int

	host     string
	port     int
	user     string
	password string

	sunnyBoyAddr *net.UDPAddr
	localAddr    *net.UDPAddr
	udpConn      *net.UDPConn
}

// NewSunnyBoyClient returns new SunnyBoyClient
func NewSunnyBoyClient(localPort int, host string, port int, user, password string) (SunnyBoyClient, error) {

	// https://github.com/AttackStrawbery/ioBroker.sma-speedwire/blob/1d213ad98a56cce3d7a683da3ff80b927ad44420/sma-speedwire.js#L212-L251

	udpNetwork := "udp"

	sunnyBoyAddr, err := net.ResolveUDPAddr(udpNetwork, fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%v", localPort))
	if err != nil {
		return nil, err
	}

	udpConn, err := net.DialUDP(udpNetwork, localAddr, sunnyBoyAddr)
	if err != nil {
		return nil, err
	}

	return &sunnyBoyClientImpl{
		localPort: localPort,

		host:     host,
		port:     port,
		user:     user,
		password: password,

		sunnyBoyAddr: sunnyBoyAddr,
		localAddr:    localAddr,
		udpConn:      udpConn,
	}, nil
}

func (sb *sunnyBoyClientImpl) writePacketHeader(packet *bytes.Buffer) {
	binary.Write(packet, binary.LittleEndian, 0x00414D53)
	binary.Write(packet, binary.LittleEndian, 0xA0020400)
	binary.Write(packet, binary.LittleEndian, 0x01000000)
	binary.Write(packet, binary.LittleEndian, 0x00000000)
	binary.Write(packet, binary.LittleEndian, 0x00)
	binary.Write(packet, binary.LittleEndian, 0x00) // Placeholder for packet length
}

func (sb *sunnyBoyClientImpl) writePacket(packet *bytes.Buffer) {

	// void writePacket(unsigned char *buf, unsigned char longwords, unsigned char ctrl, unsigned short ctrl2, unsigned short dstSUSyID, unsigned long dstSerial)

	binary.Write(packet, binary.LittleEndian, 0x65601000) // ETH_L2SIGNATURE

	// 	writeByte(buf, longwords);
	// 	writeByte(buf, ctrl);
	// 	writeShort(buf, dstSUSyID);
	// 	writeLong(buf, dstSerial);
	// 	writeShort(buf, ctrl2);
	// 	writeShort(buf, AppSUSyID);
	// 	writeLong(buf, AppSerial);
	// 	writeShort(buf, ctrl2);
	// 	writeShort(buf, 0);
	// 	writeShort(buf, 0);
	// 	writeShort(buf, pcktID | 0x8000);

}

func (sb *sunnyBoyClientImpl) writePacketTrailer(packet *bytes.Buffer) {

}

func (sb *sunnyBoyClientImpl) writePacketLength(packet *bytes.Buffer) {

}

func (sb *sunnyBoyClientImpl) Login() (success bool, err error) {

	packet := new(bytes.Buffer)

	sb.writePacketHeader(packet)

	// def tx_logon(self, password=b'0000', timeout=900):
	// if len(password) > 12:
	// 		raise ValueError
	// password += b'\x00' * (12 - len(password))
	// tag = self.gettag()

	// extra = bytearray(b'\xaa\xaa\xbb\xbb\x00\x00\x00\x00')
	// extra += bytearray(((c + 0x88) % 0xff) for c in password)
	// return self.tx_6560(self.local_addr2, self.BROADCAST2a0,
	// 										0x00010001, tag,
	// 										0x040cfffd, 7, timeout, extra)

	// https://github.com/dgibson/python-smadata2/blob/master/doc/protocol.txt

	// Packet header
	// -------------

	// Offset		Value
	// ---------------------
	// 0		0x7e
	// 1		length of packet (including header), max 0x70
	// 2		0x00
	// 3		check byte, XOR of bytes 0..2 inclusive
	// 4..9		"From" bluetooth address
	// 10..15		"To" bluetooth address
	// 16..17		Packet type (LE16)

	// 18..		Payload (format depends on packet type)

	// unsigned char pw[MAX_PWLENGTH] = {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0};

	// if (DEBUG_NORMAL) puts("logonSMAInverter()");

	// char encChar = (userGroup == UG_USER)? 0x88:0xBB;
	// //Encode password
	// unsigned int idx;
	// for (idx = 0; (password[idx] != 0) && (idx < sizeof(pw)); idx++)
	// 		pw[idx] = password[idx] + encChar;
	// for (; idx < MAX_PWLENGTH; idx++)
	// 		pw[idx] = encChar;

	// E_SBFSPOT rc = E_OK;
	// int validPcktID = 0;

	// time_t now;

	// do
	// {
	// 	pcktID++;
	// 	now = time(NULL);
	// 	writePacketHeader(pcktBuf01, addr_unknown);
	// 	writePacket(pcktBuf0EA00100, anySUSyID, anySerial);
	// 	writeLong(pcktBufFFFD040C);
	// 	writeLong(pcktBuf, userGroup);	// User / Installer
	// 	writeLong(pcktBuf00000384); // Timeout = 900sec ?
	// 	writeLong(pcktBuf, now);
	// 	writeLong(pcktBuf, 0);
	// 	writeArray(pcktBuf, pw, sizeof(pw));
	// 	writePacketTrailer(pcktBuf);
	// 	writePacketLength(pcktBuf);
	// }
	// while (!isCrcValid(pcktBuf[packetposition-3], pcktBuf[packetposition-2]));

	// ethSend(pcktBuf, inverters[inv]->IPAddress);

	// validPcktID = 0;
	// do
	// {
	// 	if ((rc = ethGetPacket()) == E_OK)
	// 	{
	// 		ethPacket *pckt = (ethPacket *)pcktBuf;
	// 		if (pcktID == (btohs(pckt->PacketID) & 0x7FFF))   // Valid Packet ID
	// 		{
	// 			validPcktID = 1;
	// 			unsigned short retcode = btohs(pckt->ErrorCode);
	// 			switch (retcode)
	// 			{
	// 				case 0: rc = E_OK; break;
	// 				case 0x0100: rc = E_INVPASSW; break;
	// 				default: rc = E_LOGONFAILED; break;
	// 			}
	// 		 }
	// 		else
	// 			if (DEBUG_HIGHEST) printf("Packet ID mismatch. Expected %d, received %d\n", pcktID, (btohs(pckt->PacketID) & 0x7FFF));
	// 	}
	// } while ((validPcktID == 0) && (rc == E_OK)); // Fix Issue 167

	// var cmdheader = "534D4100000402A00000000100";
	// var pktlength = "3A";
	// var esignature = "001060650EA0";
	// var encpasswd = "888888888888888888888888";
	// var arrayencpasswd = encpasswd.split("");
	// var cmdId = "0C04FDFF" + "07000000" + "84030000";
	// //var pass = "0000";
	// var arraypass = password.split("");
	// var timeStamp = Math.floor(Date.now() / 1000).toString(16);

	// for (var i = 0; i < password.length ; i++) {
	// 		var tmp = (parseInt((arrayencpasswd.slice(i*2,i*2+2).join('')),16) + parseInt(password.charCodeAt(i))).toString(16).slice("");
	// 		arrayencpasswd.splice(i*2,2,tmp[0],tmp[1]);
	// }
	// var cmd = cmdheader + pktlength + esignature + ByteOrderShort(anySusyId) + ByteOrderLong(anySerial);
	// cmd = cmd  + "0001" + ByteOrderShort(mySusyId) + ByteOrderLong(mySerial) + "0001" + "00000000" + decimalToHex(pktId++,4) + cmdId + timeStamp + "00000000" + arrayencpasswd.join('') + "00000000" ;
	// var cmdBytes = hex2bin(cmd);

	// socket.send(cmdBytes, 0, cmdBytes.length, PORT, HOST, function(err, bytes) {
	// 			if (err) throw err;
	// 			adapter.log.debug('login UDP message sent to ' + HOST +':'+ PORT+ ' : ' + cmd);
	//  //     client.close();
	// });
	// callBackCount++;

	return
}

func (sb *sunnyBoyClientImpl) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {

	return
	// return sb.client.ReadInputRegisters(address, quantity)
}

// /**
//  *
//  * sma-speedwire adapter
//  *
//  *
//  *  file io-package.json comments:
//  *
//  *  {
//  *      "common": {
//  *          "name":         "sma-speedwire",                  // name has to be set and has to be equal to adapters folder name and main file name excluding extension
//  *          "version":      "0.0.0",                    // use "Semantic Versioning"! see http://semver.org/
//  *          "title":        "Node.js sma-speedwire Adapter",  // Adapter title shown in User Interfaces
//  *          "authors":  [                               // Array of authord
//  *              "name <mail@sma-speedwire.com>"
//  *          ]
//  *          "desc":         "sma-speedwire adapter",          // Adapter description shown in User Interfaces. Can be a language object {de:"...",ru:"..."} or a string
//  *          "platform":     "Javascript/Node.js",       // possible values "javascript", "javascript/Node.js" - more coming
//  *          "mode":         "daemon",                   // possible values "daemon", "schedule", "subscribe"
//  *          "schedule":     "0 0 * * *"                 // cron-style schedule. Only needed if mode=schedule
//  *          "loglevel":     "info"                      // Adapters Log Level
//  *      },
//  *      "native": {                                     // the native object is available via adapter.config in your adapters code - use it for configuration
//  *          "test1": true,
//  *          "test2": 42
//  *      }
//  *  }
//  *
//  */

// /* jshint -W097 */// jshint strict:false
// /*jslint node: true */
// "use strict";

// // you have to require the utils module and call adapter func (sb *sunnyBoyClientImpl)
// var utils =    require(__dirname + '/lib/utils'); // Get common adapter utils
// var commands = {
// 	EnergyProduction: {
// 		command: littleEndianHex("54000200"),
// 		first: littleEndianHex("00260100"),
// 		last: littleEndianHex("002622FF"),
// 		label: "SPOT_ETODAY, SPOT_ETOTA"
// 	},
// 	SpotDCPower: {
// 		command: littleEndianHex("53800200"),
// 		first: littleEndianHex("00251E00"),
// 		last: littleEndianHex("00251EFF"),
// 		label: "SPOT_PDC1, SPOT_PDC2"
// 	},
// 	SpotDCVoltage: {
// 		command: littleEndianHex("53800200"),
// 		first: littleEndianHex("00451F00"),
// 		last: littleEndianHex("004521FF"),
// 		label: "SPOT_UDC1, SPOT_UDC2, SPOT_IDC1, SPOT_IDC2"
// 	},
// 	SpotACPower: {
// 		command: littleEndianHex("51000200"),
// 		first: littleEndianHex("00464000"),
// 		last: littleEndianHex("004642FF"),
// 		label: "SPOT_PAC1, SPOT_PAC2, SPOT_PAC3"
// 	},
// 	SpotACVoltage: {
// 		command: littleEndianHex("51000200"),
// 		first: littleEndianHex("00464800"),
// 		last: littleEndianHex("004655FF"),
// 		label: "POT_UAC1, SPOT_UAC2, SPOT_UAC3, SPOT_IAC1, SPOT_IAC2, SPOT_IAC3"
// 	},
// 	SpotGridFrequency: {
// 		command: littleEndianHex("51000200"),
// 		first: littleEndianHex("00465700"),
// 		last: littleEndianHex("004657FF"),
// 		label: "SPOT_FREQ"
// 	},
// 	MaxACPower: {
// 		command: littleEndianHex("51000200"),
// 		first: littleEndianHex("00411E00"),
// 		last: littleEndianHex("004120FF"),
// 		label: "INV_PACMAX1, INV_PACMAX2, INV_PACMAX3"
// 	},
// 	MaxACPower2: {
// 		command: littleEndianHex("51000200"),
// 		first: littleEndianHex("00832A00"),
// 		last: littleEndianHex("00832AFF"),
// 		label: "NV_PACMAX1_2"
// 	},
// 	SpotACTotalPower: {
// 		command: littleEndianHex("51000200"),
// 		first: littleEndianHex("00263F00"),
// 		last: littleEndianHex("00263FFF"),
// 		label: "SPOT_PACTOT"
// 	},
// 	TypeLabel: {
// 		command: littleEndianHex("58000200"),
// 		first: littleEndianHex("00821E00"),
// 		last: littleEndianHex("008220FF"),
// 		label: "INV_NAME, INV_TYPE, INV_CLASS"
// 	},
// 	SoftwareVersion: {
// 		command: littleEndianHex("58000200"),
// 		first: littleEndianHex("0082340"),
// 		last: littleEndianHex("008234FF"),
// 		label: "INV_SWVERSION"
// 	},
// 	DeviceStatus: {
// 		command: littleEndianHex("51800200"),
// 		first: littleEndianHex("00214800"),
// 		last: littleEndianHex("002148FF"),
// 		label: "INV_STATUS"
// 	},
// 	GridRelayStatus: {
// 		command: littleEndianHex("51800200"),
// 		first: littleEndianHex("00416400"),
// 		last: littleEndianHex("004164FF"),
// 		label: "INV_GRIDRELAY"
// 	},
// 	OperationTime: {
// 		command: littleEndianHex("54000200"),
// 		first: littleEndianHex("00462E00"),
// 		last: littleEndianHex("00462FFF"),
// 		label: "SPOT_OPERTM, SPOT_FEEDTM"
// 	},
// 	BatteryChargeStatus: {
// 		command: littleEndianHex("51000200"),
// 		first: littleEndianHex("00295A00"),
// 		last: littleEndianHex("00295AFF"),
// 		label: " "
// 	},
// 	BatteryInfo: {
// 		command: littleEndianHex("51000200"),
// 		first: littleEndianHex("00491E00"),
// 		last: littleEndianHex("00495DFF"),
// 		label: " "
// 	},
// 	InverterTemperature: {
// 		command: littleEndianHex("52000200"),
// 		first: littleEndianHex("00237700"),
// 		last: littleEndianHex("00618FFF"),
// 		label: " "
// 	},
// 	sbftest: {
// 		command: littleEndianHex("64020200"),
// 		first: littleEndianHex("00618C00"),
// 		last: littleEndianHex("00618FFF"),
// 		label: " "
// 	}
// };
// var cmdheader = "534D4100000402A00000000100";
// var esignature = "001060650EA0";
// var anySusyId = "FFFF";
// var anySerial = "FFFFFFFF";
// var mySusyId = "0078";
// var mySerial = "3803E8C8";
// var callBackCount = 0;
// var waitCount = 0;
// var maxWaitCount = 10;
// var dgram = require('dgram');
// var pktId = 480;

// // Configure this using the admin interface
// var PORT;
// var HOST;
// // you have to call the adapter func (sb *sunnyBoyClientImpl) and pass a options object
// // name has to be set and has to be equal to adapters folder name and main file name excluding extension
// // adapter will be restarted automatically every time as the configuration changed, e.g system.adapter.sma-speedwire.0
// var adapter = utils.adapter('sma-speedwire');

// // is called when adapter shuts down - callback has to be called under any circumstances!
// adapter.on('unload', func (sb *sunnyBoyClientImpl) (callback) {
//     try {
//         adapter.log.info('cleaned everything up...');
//         callback();
//     } catch (e) {
//         callback();
//     }
// });

// // is called if a subscribed object changes
// adapter.on('objectChange', func (sb *sunnyBoyClientImpl) (id, obj) {
//     // Warning, obj can be null if it was deleted
//     adapter.log.info('objectChange ' + id + ' ' + JSON.stringify(obj));
// });

// // is called if a subscribed state changes
// adapter.on('stateChange', func (sb *sunnyBoyClientImpl) (id, state) {
//     // Warning, state can be null if it was deleted
//     adapter.log.info('stateChange ' + id + ' ' + JSON.stringify(state));

//     // you can use the ack flag to detect if it is status (true) or command (false)
//     if (state && !state.ack) {
//         adapter.log.info('ack is not set!');
//     }
// });

// // Some message was sent to adapter instance over message box. Used by email, pushover, text2speech, ...
// adapter.on('message', func (sb *sunnyBoyClientImpl) (obj) {
//     if (typeof obj == 'object' && obj.message) {
//         if (obj.command == 'send') {
//             // e.g. send email or pushover or whatever
//             adapter.log.debug('send command');

//             // Send response in callback if required
//             if (obj.callback) adapter.sendTo(obj.from, obj.command, 'Message received', obj.callback);
//         }
//     }
// });

// // is called when databases are connected and adapter received configuration.
// // start here!
// adapter.on('ready', func (sb *sunnyBoyClientImpl) () {
//     main();
// });

// func (sb *sunnyBoyClientImpl) main() {

//     // The adapters config (in the instance object everything under the attribute "native") is accessible via
//     // adapter.config:
//     // adapter.log.info('config ip: ' + adapter.config.ip);
// 		// adapter.log.info('config port: ' + adapter.config.port);
//     // adapter.log.info('config user: ' + adapter.config.user);
// 		// adapter.log.info('config password: ' + adapter.config.password);
//     var client = dgram.createSocket('udp4');
// 		HOST = adapter.config.ip;
// 		PORT = adapter.config.port;

//     client.on('message', (msg, rinfo) => {
//       var tmp = bin2hex(msg);
//       //adapter.log.debug("server got data");
//       decodeData(tmp);
//     });

//     client.on('listening', () => {
//       const address = client.address();
//       adapter.log.debug(`server listening ${address.address}:${address.port}`);
//     });

//     //client.bind(PORT);
//     login(adapter.config.user,adapter.config.password,client);
//     sendCommand("sbftest",client);
//     sendCommand("TypeLabel",client);
//     //sendCommand("SoftwareVersion");
//     sendCommand("EnergyProduction",client);
//     sendCommand("SpotDCVoltage",client);
//     sendCommand("SpotDCPower",client);
//     sendCommand("SpotACPower",client);
//     sendCommand("SpotACVoltage",client);
//     sendCommand("SpotACTotalPower",client);
//     sendCommand("MaxACPower",client);
// 		logout(client);
// 		// Force terminate after 5min
// 		waitCallBack();
// }

// func (sb *sunnyBoyClientImpl) waitCallBack()  {
// 	//here is the trick, wait until var callbackCount is set number of callback func (sb *sunnyBoyClientImpl)s
// 	if (waitCount > maxWaitCount) {
// 		adapter.log.error("Timeout");
// 		process.exit(1);
// 	}
// 	waitCount++;
// 	if (callBackCount > 0) {
// 		adapter.log.debug("wait : "+callBackCount);
// 		setTimeout(waitCallBack, 1000);
// 		return;
// 	}
// 	process.exit(0);
// }

// // Logout
// func (sb *sunnyBoyClientImpl) logout(socket) {
// 	var cmdheader = "534D4100000402A00000000100";
// 	var pktlength = "22";
// 	var esignature = "0010606508A0";
// 	var encpasswd = "888888888888888888888888";
// 	var cmdId = "0E01FDFF" + "FFFFFFFF" + "00000000";
// 	var timeStamp = Math.floor(Date.now() / 1000).toString(16);

// 	var cmd = cmdheader + pktlength + esignature + ByteOrderShort(anySusyId) + ByteOrderLong(anySerial);

// 	cmd = cmd  + "0003" + ByteOrderShort(mySusyId) + ByteOrderLong(mySerial) + "0003" + "00000000" + decimalToHex(pktId++,4) + cmdId;
// 	var cmdBytes = hex2bin(cmd);

// 	socket.send(cmdBytes, 0, cmdBytes.length, PORT, HOST, func (sb *sunnyBoyClientImpl)(err, bytes) {
// 				if (err) throw err;
// 				adapter.log.debug('logout UDP message sent to ' + HOST +':'+ PORT + ' : ' + cmd) ;
// 	});
// 	//callBackCount++;
// }

// // Login to inverter
// func (sb *sunnyBoyClientImpl) login(user,password string,socket) {
//     cmdheader := "534D4100000402A00000000100";
//     pktlength := "3A";
//     esignature := "001060650EA0";
//     encpasswd := "888888888888888888888888";
//     arrayencpasswd := encpasswd.split("");
//     cmdId := "0C04FDFF" + "07000000" + "84030000";
//     //var pass = "0000";
//     arraypass := password.split("");
//     timeStamp := Math.floor(Date.now() / 1000).toString(16);

//     for i := 0; i < password.length ; i++ {
//         tmp := (parseInt((arrayencpasswd.slice(i*2,i*2+2).join('')),16) + parseInt(password.charCodeAt(i))).toString(16).slice("");
//         arrayencpasswd.splice(i*2,2,tmp[0],tmp[1]);
//     }
//     cmd := cmdheader + pktlength + esignature + ByteOrderShort(anySusyId) + ByteOrderLong(anySerial);
//     cmd = cmd  + "0001" + ByteOrderShort(mySusyId) + ByteOrderLong(mySerial) + "0001" + "00000000" + decimalToHex(pktId++,4) + cmdId + timeStamp + "00000000" + arrayencpasswd.join('') + "00000000" ;
//     cmdBytes := hex2bin(cmd);

// 		socket.send(cmdBytes, 0, cmdBytes.length, PORT, HOST, func (sb *sunnyBoyClientImpl)(err, bytes) {
// 		      if (err) throw err;
// 		      adapter.log.debug('login UDP message sent to ' + HOST +':'+ PORT+ ' : ' + cmd);
// 		 //     client.close();
// 		});
// 		callBackCount++;
// }

// func (sb *sunnyBoyClientImpl) sendCommand(which,socket) {
//   adapter.log.debug("called command : " + which);
//   var command = commands[which]["command"];
//   var first = commands[which]["first"];
//   var last = commands[which]["last"];
//   var ctrl2 = "0000";
//   var cmd = "534d4100";

//   // Build packet header
//   cmd += "000402a0";
//   cmd += "00000001";
//   cmd += "LEGT";  // Placeholder for Packet length
//   cmd += "00106065";		// ETH Signature

//   cmd += "09";
// 	cmd += "A0";
// 	cmd += ByteOrderShort(anySusyId);
// 	cmd += ByteOrderLong(anySerial);
// 	cmd += ctrl2;
// 	cmd += ByteOrderShort(mySusyId);
// 	cmd += ByteOrderLong(mySerial);
// 	cmd += ctrl2;
// 	cmd += "0000";
// 	cmd += "0000";
// 	cmd += decimalToHex(pktId++,4);		// Packet ID (TODO use pcktID counter)
// 	cmd += command;
// 	cmd += first;
// 	cmd += last;
// 	cmd += "00000000";

//   cmd = cmd.replace('LEGT',decimalToHex((cmd.length/2 -20),4));
//   //adapter.log.debug("cmd : "+ cmd);

//   var cmdBytes = hex2bin(cmd);

//   socket.send(cmdBytes, 0, cmdBytes.length, PORT, HOST, func (sb *sunnyBoyClientImpl)(err, bytes) {
//         if (err) throw err;
//         // adapter.log.debug('UDP message sent to ' + HOST +':'+ PORT);
//    			// client.close();
//   });
// 	callBackCount++;
// }

// func (sb *sunnyBoyClientImpl) littleEndianHex(hex) {
//   var result = "";
//   var len = hex.length;

//   while (len  > 0) {
//     len = len - 2;
//     result += hex.substr(len, 2);
//   }
//   //adapter.log.debug(hex);
//   //adapter.log.debug(result);
//   return result;
// }

// func (sb *sunnyBoyClientImpl) d2h(d) {
//     var s = (+d).toString(16);
//     if(s.length < 2) {
//         s = '0' + s;
//     }
//     return s;
// }

// func (sb *sunnyBoyClientImpl) hex2bin(hex) {
//    return new Buffer(hex,"hex");
// }

// func (sb *sunnyBoyClientImpl) bin2hex(bin) {
//   return new Buffer(bin, 'ascii').toString('hex');
// }

// func (sb *sunnyBoyClientImpl) hexToBytes(hex) {
//     var bytes = [], str;
//     for (var c = 0; c < hex.length - 1; c += 2) {
//       adapter.log.debug(hex.substr(c, 2));
//       bytes.push(parseInt(hex.substr(c, 2), 16));
//     }
//     //return String.fromCharCode.apply(String, bytes);
//     return bytes;
// }

// func (sb *sunnyBoyClientImpl) ByteOrderShort(s) {
//     var array = s.split("");
//     var output = array.slice(2, 4).join('') + array.slice(0, 2).join('');
//     return output;
// }

// // Convert Little Endian hex string to big Endian
// // Currently 32 bit fixed
// func (sb *sunnyBoyClientImpl) ByteOrderLong(s) {
//     var array = s.split("");
//     var output = array.slice(6,8).join('') + array.slice(4,6).join('') + array.slice(2,4).join('') + array.slice(0,2).join('');
//     return output;
// }

// func (sb *sunnyBoyClientImpl) decodeData(hex) {
//   adapter.log.debug("decodeData input : "+hex);

//   var loop = true;
//   var pointer = 108;
//   var long = 4; // 32 bit

// 	var cmdLength = hex.length;
// 	while (loop) {
// 		var tmp = get32Bit(ByteOrderLong(hex.substr(pointer,8)));
// 		pointer += 8;
// 		var timestamp = get32Bit(hex.substr(pointer,8));
// 		pointer +=8;
// 		//console.dir(what);
// 		var code = tmp & 0x00ffff00;
// 		var cls = tmp & 0xff;
// 		var dataType = code >> 24;
// 		var cmd = code.toString(16).toUpperCase();
// 		adapter.log.debug("cmd : " + cmd);
// 		//pointer += 16;
// 		if (cmd === "251E00") {
// 				adapter.log.debug("SPOT_PDC"+cls+" : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "821E00") {	// Device class
// 				var devClass = get32Bit(ByteOrderLong(hex.substr(pointer,8)));
// 				var type = devClass & 0x00FFFFFF;
// 				adapter.log.debug("type : "+type);
// 				pointer += 64;
// 		} else if (cmd === "821F00") { // Device class
// 				var tmp = get32Bit(ByteOrderLong(hex.substr(pointer,8))) & 0x00FFFFFF;
// 				if (tmp != 16777214 ) {
// 					adapter.setState("INV_CLASS",tmp);
// 					adapter.log.debug("INV_CLASS : "+tmp);
// 				}
// 				pointer += 64;
// 		} else if (cmd === "822000") {
// 				var tmp = get32Bit(ByteOrderLong(hex.substr(pointer,8))) & 0x00FFFFFF;
// 				if (tmp != 16777214 ) {
// 					adapter.setState("INV_TYPE",tmp);
// 					adapter.log.debug("INV_TYPE : "+tmp);
// 				}
// 				pointer += 64;
// 		} else if (cmd === "823400") {	// Software Version etc.
// 				pointer += 64;
// 		} else if (cmd === "263F00") {
// 				adapter.log.debug("SPOT_PACTOT : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_PACTOT",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "464000") {
// 				adapter.log.debug("SPOT_PAC1 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_PAC1",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "464100") {
// 				adapter.log.debug("SPOT_PAC2 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_PAC2",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "464200") {
// 				adapter.log.debug("SPOT_PAC3 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_PAC3",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "464800") {
// 				adapter.log.debug("SPOT_UAC1 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_UAC1",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "464900") {
// 				adapter.log.debug("SPOT_UAC2 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_UAC2",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "464A00") {
// 				adapter.log.debug("SPOT_UAC3 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_UAC3",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "465000") {
// 				adapter.log.debug("SPOT_IAC1 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_IAC1",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "465100") {
// 				adapter.log.debug("SPOT_IAC2 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_IAC2",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "465200") {
// 				adapter.log.debug("SPOT_IAC3 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_IAC3",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "465300") {
// 				adapter.log.debug("SPOT_IAC1_2 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_IAC1_2",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "465400") {
// 				adapter.log.debug("SPOT_IAC2_2 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_IAC2_2",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "465500") {
// 				adapter.log.debug("SPOT_IAC3_2 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_IAC3_2",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "452100") {
// 				adapter.log.debug("SPOT_IDC"+cls+" : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_IDC"+cls,get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "411E00") {
// 				adapter.log.debug("SPOT_PACMAX1 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_PACMAX1",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "411F00") {
// 				adapter.log.debug("SPOT_PACMAX2 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_PACMAX2",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "412000") {
// 				adapter.log.debug("SPOT_PACMAX3 : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_PACMAX3",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "451F00") {
// 				adapter.log.debug("SPOT_UDC"+cls+" : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				adapter.setState("SPOT_UDC"+cls,get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 				pointer += 40;
// 		} else if (cmd === "262200") {
// 			adapter.log.debug("SPOT_ETODAY : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 			adapter.setState("SPOT_ETODAY",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 			pointer += 16;
// 		} else if (cmd === "260100") {
// 			adapter.log.debug("SPOT_ETOTAL : " + get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 			adapter.setState("SPOT_ETOTAL",get32Bit(ByteOrderLong(hex.substr(pointer,8))));
// 			pointer += 16;
// 		} else {
// 			if (pointer >= cmdLength) {
// 				adapter.log.debug("End of input");
// 				loop = false;
// 			} else {
// 				adapter.log.debug("unknown cmd : "+ cmd);
// 				adapter.log.debug("cmdLength : "+cmdLength);
// 				adapter.log.debug("pointer : "+pointer);
// 				adapter.log.debug("hex : " + hex);
// 				loop = false;
// 			}
// 		}
// 	}
// 	callBackCount--;
// }

// func (sb *sunnyBoyClientImpl) get32Bit(hex) {
// 	if (hex.toUpperCase() === "80000000" || hex.toUpperCase() === "FFFFFFFF") {
// 		return 0;
// 	}
// 	return parseInt(hex,16);
// }

// func (sb *sunnyBoyClientImpl) decimalToHex(d, padding) {
//     var hex = Number(d).toString(16);
//     padding = typeof (padding) === "undefined" || padding === null ? padding = 2 : padding;

//     while (hex.length < padding) {
//         hex = "0" + hex;
//     }
//     return hex;
// }
