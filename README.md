# AudioCodesCDR

.Прием UDP:0.0.0.0:4444 по умолчанию

$ ./AudioCodesCDR -h

Usage of ./AudioCodesCDR:

* --dbase, -b, string ..... PostgreSQL database (default "cdrs")
* --debug, -d, string ..... Debug level ... 1 - log file only, 2 - console only, 3 - both, 0 - disabled (default)
* --dpasswd, -p, string ... PostgreSQL database password (default "password")
* --dserver, -s, string ... PostgreSQL database server  (default "127.0.0.1")
* --dtable, -t, string .... PostgreSQL database (default "cdrs")
* --duser, -u, string ..... PostgreSQL database user (default "user")
* --help -h ............... Show help message
* --logfile, -l, string ... Full path and name of the log file (default "/var/log/AudioCodesCDR.log")
* --netport, -n, string ... Network UDP port (default ":4444")

Таблица в базеданных создается автоматически, но надо предварительно предоставить соответствующие права пользователю.

Проверено с AudioCodes Mediant-3000 PSTN STM1\SONET Interface, firmware version 6.40A, protocol type SIP

- на AudioCodes должено в меню "System" включен "Syslog" в разделе "Syslog Settings"
- и меню "SIP Definitions" раздел "Advanced Parameters" - "CDR Report Level" = "End Call"

Двоичный файл "AudioCodesCDR" скомпилирован 4.15.0-151-generic #157-Ubuntu SMP x86\_64 GNU/Linux Mint 19 Tara


