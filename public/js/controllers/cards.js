app.controller('Cards', function cardsCtrl($rootScope, $scope, $route, $routeParams, api) {
  $scope.initCards = function() {
    $scope.players = ["ben foster", "jordi amat", "jamal blackman", "mark hudson", "costel pantilimon", "pablo zabaleta", "tomas kalas", "milos veljkovic 2", "mikel arteta", "lewis holtby", "billy jones", "stephen dobbie", "kevimcnaughton", "damien delaney", "younes kaboul", "shaquile coulthirst 2", "marouane chamakh", "joe ralls", "fabio", "gerhard tremmel", "pablo armero", "maarten stekelenburg", "marco van ginkel", "jed steer", "alan tate", "joel ward", "james milner", "nathan redmond", "gwion edwards", "mesca", "anthony pilkington", "neil taylor", "alex iwobi 2", "aaron lennon", "gaston ramirez", "paulo gazzaniga", "wayne rooney", "daniel agger", "kagishdikgacoi", "jesse joronen", "jerome thomas", "matty fryatt", "kelvin davis", "isaac hayden", "vlad chiriches", "david silva", "ashley williams 2", "thomas vermaelen", "omar rowe", "harry kane", "jo inge berget", "steve gerrard", "jack grealish", "scott parker", "saido berahino", "jermaine pennant", "kenwyne jones", "vurnon anita", "wilfried bony", "jordon mutch", "adrian mariappa", "javier garrido", "petr cech", "carlo nash", "mesuozil", "federico macheda", "paul mcshane", "oscar ustari", "jonathan de guzman", "daniel potts", "sebastien bassong", "nemanja matic", "brad friedel", "mamadou sakho", "elliott bennett", "mohamed salah", "john stones", "aaron hughes", "harrison reed", "fabio borini", "liam rosenior", "ibou touray", "pajtim kasami", "leandro bacuna", "gylfi sigurdsson", "john obi mikel", "lukasz fabianski", "olivier giroud", "lee lucas", "papiss dembcisse", "bryan oviedo", "romain amalfitano", "ryan shotton", "grant holt", "fernandinho", "robbie brady", "vincent kompany", "jan vertonghen", "reece burke", "nemanja vidic", "marc albrighton", "brad guzan", "nicolaanelka", "joel robles", "julian speroni", "libor kozak", "guillermo varela", "christian benteke", "kevin mirallas", "michel vorm", "ben nugent", "ezekiel fryers", "morgan schneiderlin", "thomas ince", "robin vapersie", "shinji kagawa", "jose enrique", "gregor zabret", "ricky van wolfswinkel", "jonas gutierrez", "joseph yobo", "chico", "yoan gouffran", "emmanuel mayuka", "ange freddy plumain", "craig bellamy", "joey o'brien", "wilson palacios", "fabian delph", "unknown", "lee cattermole", "david marshall", "jordan rossiter", "adrian", "glenn whelan", "peter crouch", "steven pienaar", "damien duff", "shay given", "cameron jerome", "josmurphy", "guy demel", "andrea dossena", "jak alnwick", "ramires", "yaya sanogo", "george moncur", "ciaran clark", "etienne capoue", "hatem ben arfa", "wayne routledge", "adam armstrong", "jose fonte 2", "alou diarra", "christian eriksen", "chris smalling", "dimitar berbatov", "per mertesacker", "aleksandar tonev", "rio ferdinand", "andy carroll", "jay rodriguez", "aiden mcgeady", "leighton baines", "mathieu flamini", "david meyler", "nathan ake", "mathieu debuchy", "jos hooiveld", "theo walcott", "ryan giggs", "ravel morrison", "samir nasri", "wes hoolahan", "simon mignolet", "seamus coleman", "gary medel", "nathaniel clyne", "david vaughan", "tikrul", "callum driver", "diego lugano", "unknown 3", "thomas davies", "hugo lloris", "muamer tankovic", "michael williamson", "jordan henderson", "charlie adam", "steven naismith", "johan elmander", "curtis davies", "emanuele giaccherini", "dejan lovren", "tomas rosicky", "michael dawson", "tommy smith 2", "daniel gabbidon", "sylvain marveaux", "zoltan gera", "daniel bachmann", "luke shaw", "thievy", "jake sinclair", "russelmartin", "josh pritchard", "danny graham", "edin dzeko", "victor moses", "ryan taylor", "matthew taylor", "johnny heitinga", "boaz myhill", "sebastian larsson", "dean moxey", "ricardo vaz te", "arouna kone", "romellukaku", "connor ogilvie 2", "steven fletcher", "conor henderson", "nikica jelavic", "frank lampard", "yacouba sylla", "roger johnson", "leon osman", "gareth mcauley", "thomas sorensen", "rafael", "valentin roberge", "eden hazard", "antonio luna", "adam lallana", "darnel situ", "james tomkins", "jason puncheon", "robert snodgrass", "joshua mceachran", "maximiliano amondarain", "victor anichebe", "danny ward 2", "john o'shea", "jamewilson", "oussama assaidi", "joe hart", "wayne hennessey", "peter whittingham", "joe cole", "lewis baker", "cala", "kim bo kyung", "davide santon", "john guidetti", "ondrej celustka", "lewis price", "emmanuel frimpong", "daniel johnson", "liam bridcutt", "jose canas", "nicky maynard", "kenneth mcevoy", "john arne riise", "nacer chadli", "jack cork", "cameron gayle", "ashley westwood", "don cowie", "marcello trotta", "claudio yacob", "sylvain distin", "loic remy", "sergio aguero", "ashley richards", "dougie samuel", "robert huth", "alex pozuelo", "adnan januzaj 2", "david cornell", "jonathan howson", "matthew lowton", "luke daniels", "aleksandakolarov", "james collins", "darron gibson", "alvaro negredo", "sam gallagher", "jernade meade", "brad smith", "ben davies 2", "william kvist", "artur boruc", "john flanagan", "laurent koscielny", "jason mccarthy", "branislav ivanovic", "philippe coutinho", "giorgos karagounis", "ignasi miquel", "andrew taylor", "asmir begovic", "roland lamah", "wilfried zaha", "hugo rodallega", "steve sidwell", "steven whittaker", "dwight gayle", "santiago vergini", "adil nabi", "phil jagielka", "jack wilshere", "abdoulaye faye", "danny whitehead", "cesar azpilicueta", "scott dann", "joe dudgeon", "ross fitzsimons", "lukas podolski", "steven n'zonzi", "erilamela", "kim kallstrom", "conor mcaleny 2", "oscar", "peter odemwingie", "matej vydra 2", "martin skrtel", "mousa dembele", "yannick sagbo", "moussa sissoko", "yaya toure", "brede hangeland", "magnus eikrem", "morgaamalfitano", "mark bunn", "bacary sagna", "stephen quinn", "garry monk", "elliot lee", "joe lewis", "joe allen", "leroy fer", "bradley jones", "juan mata", "nathan dyer", "david ngog", "reece hall johnson", "yannicbolasie 2", "david de gea", "maya yoshida", "robbert elliot", "gael kakuta", "craig dawson", "victor wanyama", "marcos alonso", "stephen ireland", "jonas olsson", "kevin phillips", "mladen petric", "ben turner", "ryashawcross", "emiliano viviano", "danny welbeck", "pablo hernandez", "tim howard", "winston reid", "mats daehli", "adel taarabt", "shola ameobi", "stephane sessegnon", "ross barkley", "charalampos mavrias", "sebastialletget", "steven reid", "mark schwarzer", "jermain defoe", "leo chambers", "kwesi appiah", "david luiz", "andy wilkinson", "darren bent", "jake livermore", "rhys healey", "glen johnson", "kieran gibbs", "youssumulumbu", "luuk de jong", "gareth barry", "paul dummett", "gary gardner", "dani pacheco", "marvin emnes", "alexander kacaniklic", "leroy lita", "mile jedinak 2", "john terry", "marcos lopes", "jonny evans", "chribrunt", "michael carrick", "ryan bertrand", "rickie lambert", "larnell cole", "joshua passley", "razvan rat", "vito mannone", "kadeem harris", "daniel ayala", "wojciech szczesny", "steven caulker", "gervinho", "paulinho", "carlos cuellar", "carlton cole", "nani", "lacina traore", "sammy ameobi", "mohamed diame", "nacho monreal", "jack colback", "keiren westwood", "gael clichy", "anders lindegaard", "steven taylor", "gediozelalem", "kyle de silva", "gary hooper", "robert koren", "modibo maiga", "conor newton", "james morrison", "unknown 2", "aaron wilbraham", "ahmed elmohamady", "karim el ahmadi", "stephen sama", "nabil bentaleb 2", "vassiriki abou diaby", "jonathan obika", "ashley young", "luis suarez", "cameron stewart", "lee camp", "scott sinclair", "erik pieters", "martin olsson", "joleon lescott", "barry bannan", "daniel sturridge", "callurobinson", "graham dorrans", "craig gardner", "james mccarthy", "chris david", "simon lappin", "adam morgan", "jussi jaaskelainen", "marc muniesa", "josmer volmy altidore", "michu", "matija nastasic", "maynor figueroa", "stuart okeefe 2", "danny rose", "jerome williams", "emmanuel adebayor", "marouane fellaini", "kyle walker", "jason beckford", "willian", "cheik tiote", "jandi donacien", "gabriel agbonlahor", "tom cleverley", "emyhuws", "alex bruce", "heurelho gomes", "santiago cazorla", "javi garcia", "jonathan walters", "david stockdale", "dan burn 2", "tommy osullivan", "aaron ramsey", "michael turner", "marco borriello", "iago aspas", "stewart downing", "kyle naughton", "jay spearing", "ignacio scocco", "dan gosling", "nathan baker", "geoff cameron", "mark noble", "ashley cole", "phil jones", "cameron brannagan", "guly", "henrique hilario", "antoninocerino", "george mccartney", "andreas weimann", "david fox", "andre schurrle", "james chester", "matthew fletcher", "adam johnson", "marko arnautovic", "steven davis", "fernando torres", "calaum jahraldo martin", "craig conway", "sone aluko", "lasse vigen christensen", "kristoffer olsson", "jimmy kebe", "jesus navas", "steve harper", "antonio valencia", "demba ba", "daniel fox", "patrick mccarthy", "joe ledley", "fabricicoloccini", "mapou yanga mbiwa", "leon britton", "dedryck boyata", "ryan fredericks", "aron gunnarsson", "ki sung yueng", "kevin theophile catherine", "gerard deulofeu", "adam smith 3", "jores okore", "el hadji ba", "kostas mitroglou", "fraizer campbell", "matthew connolly", "jonathan parr 2", "tom huddlestone", "calum chambers", "dwight tiendalli", "patrice evra", "allan mcgregor", "marc wilson", "john ruddy", "neialexander", "jonjo shelvey", "philippe senderos", "roberto soldado", "clint dempsey", "declan john", "kostamitroglou", "neil alexander", "jonjshelvey", "ron vlaar", "craig noone", "dominic ball 2", "ibra sekajja", "kevin nolan", "gary cahill", "angel rangel", "alvaro", "javier hernandez", "davrichards", "stephen henderson", "stevan jovetic", "james ward prowse", "moussa dembele", "shane long", "blair turgott", "hector bellerin", "charlie ward", "liam ridgewell", "matthew jarvis", "kieran richardson", "ryatunnicliffe", "samuel etoo", "antolin alcaraz", "kyle bartley", "fernando amorebieta", "glenn murray 2", "wes brown", "elliot grandin", "joe bennett", "raphael spiegel", "luke garbutt", "conor coady", "nicklabendtner", "adlene guedioura", "laste dombaxe 2", "matthew etherington"];
    $rootScope.cards = [];
  };

  $scope.searchCard = function() {
    var params = {
      names: [$scope.playerName.replace(" ", "-")]
    }
    api.getCard(params).then(function(data) {
      if (data && data.stats) {
        $rootScope.cards = mergeCards($scope.cards, data.stats);
        $scope.playerName = '';
      } else {
        // fix me, handle error
      }
    });
  };

  // pop form to edit an card
  $scope.editCard = function(card) {
    $scope.cardData = {
      action: 'edit',
      modalName: card.name,
      id: card.id,
      name: card.name,
      tags: card.tags,
      url: card.url
    }
  };

  function mergeCards(cards, newCards) {
    for (var i = 0; i < cards.length; i++) {
      for (var j = 0; j < newCards.length; j ++) {
        if (cards[i].nameAlias && cards[i].nameAlias === newCards[j].nameAlias) {
          newCards.splice(j, 1);
        }
      }
    }
    return cards.concat(newCards);
  }
});

app.directive('cards', function(api) {
  return {
    restrict: 'E',
    scope: {
      cards: '=',
    },
    link: function (scope, elem, attrs) {
      scope.addCard = function(name) {
        var params = {
          name: name
        }
        api.addCard(params).then(function(data) {
          console.log(data);
        });
      };

      scope.$watch('cards', function(oldVal, newVal) {
        if (scope.cards.length > 0) {
          React.renderComponent(window.Cards({players: scope.cards, addCard: scope.addCard}), elem[0]);
        }
      });
    }
  }
});
