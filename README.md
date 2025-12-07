# Bitaksi TaxiHub - Backend Microservices

Bu proje, Bitaksi Hackathon kapsamında geliştirilmiş, sürücü yönetimi ve coğrafi konumlama (geospatial) işlevlerine sahip bir backend çözümüdür. Proje, Go programlama dili kullanılarak mikroservis mimarisine uygun şekilde tasarlanmış ve Docker üzerinde çalışacak şekilde yapılandırılmıştır.

## Proje Mimarisi

Sistem, sorumlulukların ayrıldığı üç ana bileşenden oluşmaktadır:

* **API Gateway (Go + Echo):** Tüm HTTP isteklerini karşılayan giriş kapısıdır. JWT tabanlı kimlik doğrulama işlemlerini gerçekleştirir ve yetkilendirilmiş istekleri ilgili servislere yönlendirir (Reverse Proxy).
* **Driver Service (Go + Standart Kütüphane):** İş mantığının yürütüldüğü servistir. MongoDB ile veri alışverişini sağlar ve Haversine formülü kullanarak koordinat bazlı mesafe hesaplamaları yapar.
* **MongoDB:** Sürücü profillerini ve enlem/boylam verilerini saklayan veritabanı katmanıdır.

## Kullanılan Teknolojiler

* **Dil:** Go (Golang) 1.24
* **Web Framework:** Echo (Gateway servisi için)
* **Veritabanı:** MongoDB
* **Sanallaştırma:** Docker & Docker Compose
* **Dokümantasyon:** Swagger (OpenAPI)
* **Güvenlik:** JWT (JSON Web Token)

## Kurulum ve Çalıştırma

Proje, yerel geliştirme ortamında `docker-compose` kullanılarak tek komutla ayağa kaldırılabilir.

### Gereksinimler

* Docker
* Docker Desktop (veya Daemon)

### Kurulum Adımları

1.  Depoyu yerel makinenize klonlayın:
    ```bash
    git clone [https://github.com/KULLANICI_ADINIZ/bitaksi-hackathon.git](https://github.com/KULLANICI_ADINIZ/bitaksi-hackathon.git)
    cd bitaksi-hackathon
    ```

2.  Servisleri derleyin ve başlatın:
    ```bash
    docker-compose up --build
    ```

3.  Konteynerler ayağa kalktığında servisler aşağıdaki portlarda çalışacaktır:

    * **API Gateway:** http://localhost:8000
    * **Swagger UI (Dokümantasyon):** http://localhost:8000/swagger/index.html
    * **Mongo Express (Veritabanı Arayüzü):** http://localhost:8081

## API Kullanımı

Sistemdeki endpoint'ler JWT koruması altındadır. İstek yapmadan önce token alınması gerekmektedir.

### 1. Kimlik Doğrulama (Token Alma)
Gateway üzerinden sisteme giriş yaparak Bearer token alınır.

* **Endpoint:** `POST /login`
* **Body:**
    ```json
    {
      "username": "admin",
      "password": "password123"
    }
    ```

### 2. Sürücü İşlemleri
Elde edilen token, header bilgisine `Authorization: Bearer <TOKEN>` formatında eklenerek istek yapılır.

* **Sürücü Listeleme:** `GET /drivers`
* **Yakındaki Sürücüleri Bulma:** `GET /drivers/nearby?lat=41.0&lon=29.0&taxiType=yellow`

---
*Bitaksi Hackathon süreci için geliştirilmiştir.*