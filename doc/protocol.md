# Protocol

## login_packet

### DrComHeader 0-3

0x03 0x01 0x00 len(username)+20->1B

### PasswordMD5 4-19

MD5(0x03 0x01 seed(4B) password(var))->16B

### username 20-55

username(var)

### ControlCheckStatus 56

ControlCheckStatus

### AdapterNum 57

AdapterNum

### MacAddrXORPasswordMD5 58-63

MAC^PasswordMD5

### PasswordMD5_2 64-79

MD5(0x01 password(var) seed(4B))->16B

### HostIpNum 80

0x01

### HostIpList 81-96

HostIP(4B), 0x00*12

### HalfMD5 97-104

MD5(login_packet(0-96),0x14,0x00,0x07,0x0b)->first half, 8b

### DogFlag 105

IPDOG

### host_name 110-141

host_name(var)

### PrimaryDNS 142-145

PrimaryDNS(4B)

### DHCP 146-149

DHCP(4B)

### OSVersionInfoSize 162-165

0x94 0x00 0x00 0x00

### MajorVersion 166-169

0x10 0x00 0x00 0x00

### MinorVersion 170-173

0x00 0x00 0x00 0x00

### BuildNumber 174-177

0x00 0x28 0x00 0x00

### PlatformID 178-181

0x02 0x00 0x00 0x00

### AuthVersion 310-311

AUTH_VERSION(2B)

*Non-ROR* version

### Code 312

0x02

### Len 313

0x0c

### CRC 314-317

crc(login_packet(314B),0x01,0x26,0x07,0x11,mac(6B))->4B

```c
    sum = 1234;
    uint64_t ret = 0;
    for (int i = 0; i < counter + 14; i += 4) {
        ret = 0;
        // reverse unsigned char array[4]
        for (int j = 4; j > 0; j--) {
            ret = ret * 256 + (int)checksum2_str[i + j - 1];
        }
        sum ^= ret;
    }
    sum = (1968 * sum) & 0xffffffff;
    for (int j = 0; j < 4; j++) {
        checksum2[j] = (unsigned char)(sum >> (j * 8) & 0xff);
    }
```

### AdapterAddress 320-325

MAC(6B)

### BroadcastMode 326

0xe9

### unknown 327

0x13

## Keepalive1

*no* keepalive1_mod

### 