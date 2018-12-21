Poniżej zostały opisane czynności niezbędne do wykonania w celu uruchomienia
środowiska, w którym będzą działać programy `counter` oraz `twittervotes`:

Każdy z poniższych programów należy uruchomić w osobnym oknie terminala:

    mongod

    nsqlookupd

    nsqd --lookupd-tcp-address=localhost:4160

Następnie należy uruchomić programy `counter` oraz `twittervotes`; przy czym każdy
z nich także powinien być uruchomiony w osobnym oknie terminala.

Uwaga:
  Zgodnie z informacjami zamieszczonymi w tekście książki, uruchomienie programu twittervotes.exe
  wymaga ustawienia kilku zmiennych środowiskowych. W systemie Linux można to zrobić wykonując 
  plik wsadowy setup.sh. Natomiast w przypadku systemu Windows można skorzystać z pliku 
  twittervotes.bat, który ustawi niezbędne zmienne środowiskowe i uruchomi program twittervotes.exe. 

Uruchomić wiersz poleceń bazy danych MondoDB:

    mongo

Utworzyć nową kolekcję o nazwie `polls` w bazie danych `ballots`:

    use ballots
    db.polls.insert({"title":"Ankieta testowa","options":["zadowolony","smutny","porażka","zwycięstwo"]})

Praktyka pokazuje, że znacznie więcej wyników można uzyskać wyszukując na Twitterze słowa w języku 
angielskim, dlatego zamiast powyższego polecenia można zacząć do tego przedstawionego poniżej:

    db.polls.insert({title:"Ankieta testowa",options:["happy","sad,"win","fail"]})

After a while, see the results by printing the polls:
Po chwili można już sprawdzić wyniki wyświetlając zawartość kolekcji `polls`:

    db.polls.find().pretty()

Po przygotowaniu bazy danych danych należy uruchomić program zliczający głosy z katalogu counter
oraz program analizujący wiadomości publikowane na Twitterze (korzystając z pliku wsadowego). 

Następnie należy uruchomić serwer API wykonując program main z katalogu api.

I w końcu można uruchomić witrynę WWW której serwer jest dostępny w katalogu web. Po wykonaniu
tych wszystkich czynności uzyskamy działające rozwiązanie, które analizuje Twittera i zapisuje 
wyniki w bazie MongoDB i udostępnia witrynę za pośrednictwem której można tworzyć ankiety 
i wyświetlać ich wyniki.