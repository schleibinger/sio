/*

1) Original: Copyright (c) 2005-2008 Dustin Sallings <dustin@spy.net>. 

2) Mods: Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim GmbH
<info@schleibinger.com>. Blame: Jan Mercl.

All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

*/

#include "_cgo_export.h"

void initRates() {
#ifdef B0
    addRate(0, B0);
#endif
#ifdef B50
    addRate(50, B50);
#endif
#ifdef B75
    addRate(75, B75);
#endif
#ifdef B110
    addRate(110, B110);
#endif
#ifdef B134
    addRate(134, B134);
#endif
#ifdef B150
    addRate(150, B150);
#endif
#ifdef B200
    addRate(200, B200);
#endif
#ifdef B300
    addRate(300, B300);
#endif
#ifdef B600
    addRate(600, B600);
#endif
#ifdef B1200
    addRate(1200, B1200);
#endif
#ifdef B1800
    addRate(1800, B1800);
#endif
#ifdef B2400
    addRate(2400, B2400);
#endif
#ifdef B4800
    addRate(4800, B4800);
#endif
#ifdef B7200
    addRate(7200, B7200);
#endif
#ifdef B9600
    addRate(9600, B9600);
#endif
#ifdef B14400
    addRate(14400, B14400);
#endif
#ifdef B19200
    addRate(19200, B19200);
#endif
#ifdef B28800
    addRate(28800, B28800);
#endif
#ifdef B38400
    addRate(38400, B38400);
#endif
#ifdef B57600
    addRate(57600, B57600);
#endif
#ifdef B76800
    addRate(76800, B76800);
#endif
#ifdef B115200
    addRate(115200, B115200);
#endif
#ifdef B230400
    addRate(230400, B230400);
#endif
#ifdef B460800
    addRate(460800, B460800);
#endif
#ifdef B500000
    addRate(500000, B500000);
#endif
#ifdef B576000
    addRate(576000, B576000);
#endif
#ifdef B921600
    addRate(921600, B921600);
#endif
#ifdef B1000000
    addRate(1000000, B1000000);
#endif
#ifdef B1152000
    addRate(1152000, B1152000);
#endif
#ifdef B1500000
    addRate(1500000, B1500000);
#endif
#ifdef B2000000
    addRate(2000000, B2000000);
#endif
#ifdef B2500000
    addRate(2500000, B2500000);
#endif
#ifdef B3000000
    addRate(3000000, B3000000);
#endif
#ifdef B3500000
    addRate(3500000, B3500000);
#endif
#ifdef B4000000
    addRate(4000000, B4000000);
#endif
}
