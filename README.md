# Bitaksi TaxiHub - Full Stack Microservices

Bu proje, Bitaksi GO Bootcampi süreci kapsamında geliştirilmiştir. Sürücü yönetimi, coğrafi konumlama ve harita üzerinde görselleştirme işlevlerine sahip, mikroservis mimarisine uygun bir full-stack uygulamadır.

Backend tarafı Go ile, Frontend tarafı React ile geliştirilmiş olup, tüm sistem Docker üzerinde çalışacak şekilde yapılandırılmıştır.

## Proje Mimarisi

Sistem dört ana bileşenden oluşmaktadır:

* **Frontend (React + Vite + Leaflet):** Son kullanıcının harita üzerinde taksileri görüntülediği, araç tipine (Sarı/Siyah) göre filtreleme yapabildiği arayüz.
* **API Gateway (Go + Echo):** Sistemin giriş kapısıdır. HTTP isteklerini karşılar, JWT tabanlı kimlik doğrulamayı yapar ve istekleri ilgili servislere yönlendirir (Reverse Proxy).
* **Driver Service (Go):** Temel iş mantığının bulunduğu servistir. MongoDB ile iletişim kurar ve Haversine formülü ile koordinat bazlı mesafe hesaplamalarını gerçekleştirir.
* **MongoDB:** Sürücü verilerini ve konum bilgilerini saklayan veritabanıdır.

## Teknik Altyapı

**Backend**
* Go (Golang) 1.24
* Echo Framework (Gateway)
* MongoDB & Mongo Driver
* JWT (JSON Web Token)
* Swagger (OpenAPI)

**Frontend**
* React.js
* Leaflet (Harita Entegrasyonu)
* Tailwind CSS
* Axios

**DevOps**
* Docker & Docker Compose

## Kurulum ve Çalıştırma

Proje, docker-compose yapılandırması sayesinde tek komutla çalıştırılabilir. Bilgisayarınızda Docker Desktop'ın yüklü olması yeterlidir.

1.  **Projeyi Klonlayın:**
    ```bash
    git clone [https://github.com/eneszeyt/bitaksi-TaxiHub.git](https://github.com/eneszeyt/bitaksi-TaxiHub.git)
    cd bitaksi-TaxiHub
    ```

2.  **Servisleri Başlatın:**
    ```bash
    docker-compose up --build
    ```
    Bu komut veritabanını, backend servislerini ve frontend uygulamasını sırasıyla hazırlar ve başlatır.

3.  **Erişim Adresleri:**
    Konteynerler ayağa kalktığında aşağıdaki adreslerden erişim sağlayabilirsiniz:

    * **Frontend Arayüzü:** http://localhost:5173
    * **API Gateway:** http://localhost:8000
    * **Swagger Dokümantasyonu:** http://localhost:8000/swagger/index.html
    * **Mongo Express:** http://localhost:8081

## Kullanım Detayları

**1. Arayüz Girişi**
Frontend arayüzüne (http://localhost:5173) aşağıdaki varsayılan bilgilerle giriş yapabilirsiniz:
* Kullanıcı Adı: `admin`
* Şifre: `password123`

**2. Harita Özellikleri**
Giriş yaptıktan sonra harita üzerinde kayıtlı sürücüleri görebilir, sağ üstteki butonlar ile sarı veya siyah taksi filtresi uygulayabilirsiniz. Sürücü pinlerine tıklayarak detayları görüntüleyebilirsiniz.

**3. API Kullanımı**
API'yi harici olarak kullanmak isterseniz `POST /login` endpoint'inden token almanız ve diğer isteklere `Authorization: Bearer <TOKEN>` başlığını eklemeniz gerekmektedir.

---
Bitaksi TaxiHub projesidir.