# bankProject
An XML to DB, DB to XML project, which was builded on GO.

Main file is bankProject.go, there are some settings in this file, which affect everything else in the package.
Right at the beginning of the main() function you will see "settings blog" with two settings:

"ConnectionDB" is string for mysql connection handler, in this string you can choose your user, password, port and database to connect. 

"settings" string is responsible for switching modes: Import and Retrieve, but the whole functioanlity of this feature is now deactivated due to me not realising how to do a "Retrieve" part properly("I've made it really really close to how it should look in the end, but time's up). In this version only "Import" is available. Import takes data from "bank.Credentials.xml" and imports it in the database, checking and creating tables by his own. Retrieve should take data from database and convert it into the XML file.

----------------------
Database creation

This application requires mysql database server installed on your pc, which can be downloaded here: https://dev.mysql.com/downloads/installer/

Once you have installed it, created user/password and chosen a port to start a mysql server, log into your account and create database named "bankProject". OR you can name it how you want, but be sure to change the name of database in the "ConnectionDB" string which is located in the bankProject.go file.

To operate with databases GoLang requires drivers, which are needed to be installed from: https://github.com/go-sql-driver/mysql

WARNING: all .go files in the repository are one package, so they are needed to start simultaneously(e.g. from CMD the command will be "go run bankProject.go dbAccount.go dbCard.go dbClient.go" 

When the database is created, and the name of your database is similar to the name of database written in "ConnectionDB", all drivers are installed, you can start the script with the command which was mentiond earlier. 

The script will automatically connect to the database you have mentioned in the "ConnectionDB" string and check the tables which are need for its work. It will create all missing tables and fill them with data from "bankCredentials.xml"
 
