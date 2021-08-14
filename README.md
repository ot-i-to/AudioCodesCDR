# AudioCodesCDR

Прием UDP:0.0.0.0:4444 по умолчанию

\$ ./AudioCodesCDR -h Usage of ./AudioCodesCDR: -b, --dbase string PostgreSQL database (default "cdrs") -d, --debug string Debug level ... 1 - log file only, 2 - console only, 3 - both, 0 - disabled (default) (default "0") -p, --dpasswd string PostgreSQL database password (default "password") -s, --dserver string PostgreSQL database server (default "127.0.0.1") -t, --dtable string PostgreSQL database (default "cdrs") -u, --duser string PostgreSQL database user (default "user") -h, --help Show help message -l, --logfile string Full path and name of the log file (default "/var/log/AudioCodesCDR.log") -n, --netport string Network UDP port (default ":4444")

проверено с AudioCodes Mediant-3000 PSTN STM1Interface, firmware version 6.40A, protocol type SIP ... на AudioCodes должено быть в меню "System" включен "Syslog" в разделе "Syslog Settings" ... и меню "SIP Definitions" раздел "Advanced Parameters" - "CDR Report Level" = "End Call"

скомпилированно 4.15.0-151-generic \#157-Ubuntu SMP Fri Jul 9 23:07:57 UTC 2021 x86\_64 x86\_64 x86\_64 GNU/Linux Distributor ID: LinuxMint Description: Linux Mint 19 Tara Release: 19 Codename: tara
