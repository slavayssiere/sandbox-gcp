function onOpen(e) {
  // Add a custom menu to the spreadsheet.
  SpreadsheetApp.getUi() // Or DocumentApp, SlidesApp, or FormApp.
      .createMenu('Get data')
      .addItem('Get users top 10', 'runTop10')
      .addSeparator()
      .addItem('Get last 10', 'searchLast10Tweet')
      .addToUi();
}

function runTop10() {

  var url = 'http://public.gcp-wescale.slavayssiere.fr/aggregator/top10';
  var options = {
    'method' : 'get'
  };
  var response = UrlFetchApp.fetch(url, options)
  var tags = JSON.parse(response.getContentText());
  
  var spreadsheet = SpreadsheetApp.getActiveSpreadsheet()
  var sheet = spreadsheet.getActiveSheet()
  
  var column = 2
  tags.forEach(function(tag){
    Logger.log(tag['tag']);
    sheet.getRange(1, column).setValue(tag['tag'])
    var users = tag['users']
    var row = 3
    users.forEach(function(user){
      sheet.getRange(row, column).setValue(user['count'])
      sheet.getRange(row, column+1).setValue(user['user'])
      row++
    });
    column=column+2
  });
}
