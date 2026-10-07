package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Uç ajanlardan (ModbusTrap) gelecek veri paketinin şablonu
type ThreatAlert struct {
	AgentID    string `json:"agent_id"`
	AttackerIP string `json:"attacker_ip"`
	AttackType string `json:"attack_type"`
	RawHex     string `json:"raw_hex"`
	Timestamp  string `json:"timestamp"`
}

var (
	// Aynı anda yüzlerce ajandan veri gelirse çakışmayı önlemek için bellek kilidi
	mutex sync.Mutex
)

func main() {
	// Ajanların rapor göndereceği API Uç Noktası (Endpoint)
	http.HandleFunc("/api/alert", handleAlert)

	fmt.Println("[+] SCADA Merkezi Komuta Kontrol (SOC/C2) Sunucusu Başlatıldı")
	fmt.Println("[+] 8080 Portu Üzerinden Tesis Ajanlarından (ModbusTrap) Veri Bekleniyor...")
	fmt.Println("========================================================================")

	// C2 sunucusunu 8080 portundan ayağa kaldırıyoruz
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handleAlert(w http.ResponseWriter, r *http.Request) {
	// Ajanlar veriyi POST metoduyla göndermek zorunda
	if r.Method != http.MethodPost {
		http.Error(w, "Erişim Reddedildi: Sadece POST istekleri kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	var alert ThreatAlert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, "Bozuk Veri Paketi", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	// İleride bu verileri bir veritabanına (PostgreSQL/MongoDB) yazabiliriz
	// Şimdilik doğrudan SOC terminaline basıyoruz
	fmt.Printf("\n[🔥 SOC ALARMI] Tesisin Uç Noktasından Yeni Tehdit Raporlandı!\n")
	fmt.Printf(" ├─ Raporlayan Ajan: %s\n", alert.AgentID)
	fmt.Printf(" ├─ Saldırgan IP   : %s\n", alert.AttackerIP)
	fmt.Printf(" ├─ Tehdit Tipi    : %s\n", alert.AttackType)
	fmt.Printf(" ├─ Ham Veri (Hex) : %s\n", alert.RawHex)
	fmt.Printf(" └─ Zaman          : %s\n", alert.Timestamp)
	fmt.Println("------------------------------------------------------------------------")
	mutex.Unlock()

	// Ajana "Rapor alındı" bilgisini dönüyoruz
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("C2: Alarm başarıyla merkeze iletildi."))
}
