# Node.js imajını kullan (Senin versiyonuna uygun)
FROM node:18-alpine

# Çalışma dizinini ayarla
WORKDIR /app

# Paket dosyalarını kopyala
COPY package.json package-lock.json ./

# Kütüphaneleri yükle
RUN npm install

# Proje dosyalarını kopyala
COPY . .

# Vite'ın dışarıdan erişilebilmesi için portu aç
EXPOSE 5173

# Uygulamayı başlat
CMD ["npm", "run", "dev", "--", "--host"]