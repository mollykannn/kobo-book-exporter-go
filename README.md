# kobo-book-exporter-go

Export Kobo Book Highlight (.md) & Book List (.csv, .json & .md).

Web版本在這裡 => [https://mollykannn.github.io/kobo-book-exporter](https://mollykannn.github.io/kobo-book-exporter)


Retrieved from [Kobo Exporter: 匯出 Kobo 電子書的書籍清單與註記資料 (劃線與筆記) | Vixual](http://www.vixual.net/blog/archives/117)

<br/>

## 使用方法

### 1) 下載此資料夾

<br/>

### 2) 拿取koboReader.sqlite

<br/>

有兩種方法可以取得KoboBook裡的sqlite :
  1) 在Kobo Reader裡尋找

      將Kobo Reader連上電腦後，進入資料夾，會看到一個.Kobo的資料夾。進去然後將KoboReader.sqlite複製到這個資料夾裡就可以了。

      (.Kobo資料夾是隱藏的。在MacOS裡按cmd+Shift+.就可以看到資料夾了)

  <br/>

  2) 在Kobo Desktop Application裡尋找

      2.1. Windows

        進到 C:\Users\{用戶名稱}\AppData\Local\Kobo\Kobo Desktop Edition\，將 Kobo.sqlite 複製到這個資料夾，然後改名為KoboReader.sqlite。

      2.2. MAC OS
      
        進到 /Users/{用戶名稱}/Library/Application Support/Kobo/Kobo Desktop Edition，將 Kobo.sqlite 複製到這個資料夾，然後改名為KoboReader.sqlite。
        (PS: Library資料夾預設是隱藏的，需要按cmd+Shift+.才能看見)

<br/>

### 3) 運行程式

<br/>

Windows - 運行koboBookExport_win.exe

MacOS - 開啟terminal，進到此資料夾路徑，然後打上 ./koboBookExport_mac

Linux - 開啟terminal，進到此資料夾路徑，然後打上 ./koboBookExport_linux
