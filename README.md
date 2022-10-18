# aufgabe04b

Implementieren Sie mit Hilfe von gRPC und Protobuffern ein einfaches
Server/Client Paar.
Der Client kann vom Server die Inhalte einer Textdatei anfordern. Diese
Inhalte werden dem Clienten dann zeilenweise als Datenstrom gesendet. Stellen
Sie sicher, dass der Server bei einem Verbindungsabbruch das Senden der Daten
abbricht.
