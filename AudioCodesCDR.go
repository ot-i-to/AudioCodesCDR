package main

// https://docs.microsoft.com/ru-ru/azure/postgresql/connect-go

import (
	//"encoding/hex"
	"bufio"
	"database/sql"
	"fmt"

	"log"
	"net"
	"os"
	"strings"

	"github.com/spf13/pflag"

	_ "github.com/lib/pq"
)

var (
	// Initialize connection constants.
	HOST     = "tdb01.almatel.msk.ru"
	DATABASE = "cdrs"
	TABLE    = "ac_cdr"
	USER     = "audiocodescdr"
	PASSWORD = "Pdjyrb!"
	// Network variables
	PORT     = ":4444"
	PROTOCOL = "udp"
	// Basic variables
	LOGFILE  = "/var/log/AudioCodesCDR.log"
	DEBUG    = "0"
	showHelp bool
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func standardizeSpaces(s string) string {
	//return strings.Join(strings.Fields(s), " ")
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

//////////////////////////////////////////////////////////////////////

func main() {

	pflag.StringVarP(&DEBUG, "debug", "d", "0", "Debug level ... 1 - log file only, 2 - console only, 3 - both, 0 - disabled (default)")
	pflag.StringVarP(&PORT, "netport", "n", ":4444", "Network UDP port")
	pflag.StringVarP(&LOGFILE, "logfile", "l", "/var/log/AudioCodesCDR.log", "Full path and name of the log file")
	pflag.StringVarP(&HOST, "dserver", "s", "127.0.0.1", "PostgreSQL database server ")
	pflag.StringVarP(&DATABASE, "dbase", "b", "cdrs", "PostgreSQL database")
	pflag.StringVarP(&TABLE, "dtable", "t", "cdrs", "PostgreSQL database")
	pflag.StringVarP(&USER, "duser", "u", "user", "PostgreSQL database user")
	pflag.StringVarP(&PASSWORD, "dpasswd", "p", "password", "PostgreSQL database password")

	pflag.BoolVarP(&showHelp, "help", "h", false,
		"Show help message")
	pflag.Parse()
	if showHelp {
		pflag.Usage()
		return
	}

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	//db, err := sql.Open("postgres", "postgres://user:pass@localhost/bookstore")
	//if err != nil {
	//    log.Fatal(err)
	//}

	// Initialize connection string.
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE)

	// Initialize connection object.
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully created connection to database")

	defer db.Close()

	//t := time.Now()
	//year := t.Year()
	//fmt.Println(year)

	// проверка наличие наблицы в бд
	var tblcount int
	sql_statement := "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' and table_name like '" + TABLE + "';"
	rows, err := db.Query(sql_statement)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		switch err := rows.Scan(&tblcount); err {
		case sql.ErrNoRows:
			log.Println("No rows were returned")
		case nil:
			switch tblcount {
			case 1:
				// при наличчи, работаем дальше
				log.Printf("Ok count Table = %d\n", tblcount)
			case 0:
				// при отсутсвии, создать
				log.Printf("Not count Table = %d\n", tblcount)
				// Create table.
				sql_cteatetable := "CREATE TABLE " + TABLE + " (id bigserial NOT NULL, PRIMARY KEY (id), INhost text NULL, Cid  text NULL, SessionId text NULL, Trunk integer NOT NULL, BChan integer NOT NULL, ConId integer NOT NULL, TG integer NOT NULL, EPTyp text NOT NULL, Orig text NOT NULL, SourceIp text NOT NULL, DestIp text NOT NULL, SrcTON integer NOT NULL, SrcNPI integer NOT NULL, SrcPhoneNum text NOT NULL, SrcNumBeforeMap text NOT NULL, DstTON integer NOT NULL, DstNPI integer NOT NULL, DstPhoneNum text NOT NULL, DstNumBeforeMap text NOT NULL, Durat integer NOT NULL, Coder text NOT NULL, Intrv integer NOT NULL, RtpIp text NOT NULL, Port text NOT NULL, TrmSd text NOT NULL, TrmReason text NOT NULL, Fax integer NOT NULL, InPackets integer NOT NULL, OutPackets integer NOT NULL, PackLoss integer NOT NULL, RemotePackLoss integer NOT NULL, SIPCallId text NOT NULL, SetupTime timestamp NOT NULL, ConnectTime timestamp NOT NULL, ReleaseTime timestamp NOT NULL, RTPdelay integer NOT NULL, RTPjitter integer NOT NULL, RTPssrc text NULL, RemoteRTPssrc text NULL, RedirectReason integer NOT NULL, TON integer NOT NULL, NPI integer NOT NULL, RedirectPhonNum text NOT NULL, MeteringPulses integer NOT NULL, SrcHost text NOT NULL, SrcHostBeforeMap text NOT NULL, DstHost text NOT NULL, DstHostBeforeMap text NOT NULL, IPG_description text NOT NULL, LocalRtpIp text NOT NULL, LocalRtpPort text NOT NULL, Amount text NULL, Mult text NULL, TrmReasonCategory text NOT NULL, RedirectNumBeforeMap text NULL, SrdId_name text NULL, SIPInterfaceId integer NOT NULL, ProxySetId integer NOT NULL, IpProfileId_name text NOT NULL, MediaRealmId_name text NOT NULL, SigTransportType text NOT NULL, TxRTPIPDiffServ integer NOT NULL, TxSigIPDiffServ integer NOT NULL, LocalRFactor integer NOT NULL, RemoteRFactor integer NOT NULL, LocalMosCQ integer NOT NULL, RemoteMosCQ integer NOT NULL, SigSourcePort text NOT NULL, SigDestPort text NOT NULL, MediaType text NOT NULL, GWReportType text NULL);"
				_, err = db.Exec(sql_cteatetable)
				checkError(err)
				log.Println("Finished creating table")
			default:
				// при кол-ве болше 1 или отрезательном результате, выходим из программы
				log.Printf("Count Table = %d\n", tblcount)
				log.Panic("ERROR count Table returned")
			}
		default:
			checkError(err)
		}
	}

	//Build the address ... "Wrong Address"
	udpAddr, err := net.ResolveUDPAddr(PROTOCOL, PORT)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reading " + PROTOCOL + " from " + udpAddr.String())

	//Create the connection
	udpConn, err := net.ListenUDP(PROTOCOL, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer udpConn.Close()

	//Keep calling this function
	data := make([]byte, 65535)
	var outstring string
	for {
		_, remoteAddr, err := udpConn.ReadFromUDP(data)
		checkError(err)

		var remAddr string
		remAddr = remoteAddr.String()

		// прием данных
		outmess, err := bufio.NewReader(udpConn).ReadString('\n')
		if err == nil {
			// разделаем строку на фрагменты по разделителю
			words := strings.Split(standardizeSpaces(outmess), "|")
			if len(words) == 71 {
				if words[0] == "<142>" && words[1] == "CALL_END " {
					if len(words[34]) < 4 {
						words[34] = words[33]
					}
					var tobd string
					for idx := 2; idx < len(words); idx++ {
						words[idx] = strings.ReplaceAll(standardizeSpaces(words[idx]), " (", "_")
						words[idx] = strings.ReplaceAll(standardizeSpaces(words[idx]), ")", "")
						tobd = tobd + "', '" + standardizeSpaces(words[idx])
					}
					outstring = "'" + standardizeSpaces(remAddr) + tobd + "', '" + standardizeSpaces(words[1]) + "'"
				}
			}
		}

		// Insert some data into table.
		sql_statement := "INSERT INTO ac_cdr (inhost, cid, sessionid, trunk, bchan, conid, tg, eptyp, orig, sourceip, destip, srcton, srcnpi, srcphonenum, srcnumbeforemap, dstton, dstnpi, dstphonenum, dstnumbeforemap, durat, coder, intrv, rtpip, port, trmsd, trmreason, fax, inpackets, outpackets, packloss, remotepackloss, sipcallid, setuptime, connecttime, releasetime, rtpdelay, rtpjitter, rtpssrc, remotertpssrc, redirectreason, ton, npi, redirectphonnum, meteringpulses, srchost, srchostbeforemap, dsthost, dsthostbeforemap, ipg_description, localrtpip, localrtpport, amount, mult, trmreasoncategory, redirectnumbeforemap, srdid_name, sipinterfaceid, proxysetid, ipprofileid_name, mediarealmid_name, sigtransporttype, txrtpipdiffserv, txsigipdiffserv, localrfactor, remoterfactor, localmoscq, remotemoscq, sigsourceport, sigdestport, mediatype, gwreporttype) VALUES (" + outstring + ");"
		if DEBUG == "1" {
			log.Printf("SQL --> %s\n", sql_statement)
		}
		if DEBUG == "2" {
			fmt.Printf("SQL --> %s\n", sql_statement)
		}
		if DEBUG == "3" {
			log.Printf("SQL --> %s\n", sql_statement)
			fmt.Printf("SQL --> %s\n", sql_statement)
		}
		_, err = db.Exec(sql_statement)
		//_, err = db.Exec(sql_statement, out)
		if err != nil {
			log.Printf("ERROR SQL --> %s\n", sql_statement)
		}
		checkError(err)
	}
}

///////////////////////////////////////////////////////////////////
