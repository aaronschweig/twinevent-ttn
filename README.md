# Twinevent - TTN Connector

Der TTN-Connector stellt sicher, dass ein SmartCube, die entsprechende Konfiguration erhält. Dafür werden unter anderem Device-Informationen von TTN benötigt.

## Ablauf

1. Ein SmartCube wird eingeschaltet und besitzt eine stabile Internetverbindung.
2. Der SmartCube published eine Nachricht auf ein Topic einen MQTT-Brokers. Dabei verwendet er folgendes Format: `registration/<MAC-Adresse des Cubes>`
3. Der TTN-Connector subscribed auf das Pattern `registration/+`, um alle Mac-Adressen abzufangen
4. Der TTN-Connector überpüft, ob ein korrespondierendes Device in TTN vorhanden ist
   1. **Falls nicht**, erstellt er ein Device für die mitgesendete MAC-Adresse und die Informationen werden an den SmartCube zurückgesendet
   2. **Falls ja\*** werden die benötigten Informationen an den SmartCube zurückgesendet
5. Der TTN-Connector überprüft als nächstes, ob das Gerät bereits in Eclipse-Ditto als DT angelegt ist. Auch hier wird im Falle des nichtvorhandenseins ein neues Gerät erstellt

## Konfiguration

| Environment-Variable | Information                                                                                                |
| -------------------- | ---------------------------------------------------------------------------------------------------------- |
| `TTN_ACCESS_KEY`     | Acces-Key für die Applikation in TTN. Dieser wird für die Verbindung zu TTN benötigt                       |
| `TTN_APP_ID`         | Application-ID aus TTN. Definiert die Applikation, in der Devices angelegt werden                          |
| `MQTT_BROKER`        | Adresse **inklusive Port** des zur Registrierung verwendeten MQTT-Brokers. Beispiel: `mq.jreiwald.de:1883` |
| `MQTT_USER`          | Nutzername für den Zugriff auf den Broker                                                                  |
| `MQTT_PASSWORD`      | Passwort für den Zugriff auf den Broker                                                                    |
