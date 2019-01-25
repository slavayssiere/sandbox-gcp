function listDataset() {
  
  var spreadsheet = SpreadsheetApp.getActiveSpreadsheet()
  var sheet = spreadsheet.getSheetByName('Top10')
  var listDs = BigQuery.Datasets.list("slavayssiere-sandbox").datasets
  var i = 0
  listDs.forEach(function(element) {
    Logger.log(element)
    sheet.getRange(2 + i, 1).setValue(element.id);
    i++
  });
 
}




function searchLast10Tweet(){
  var spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
  var sheet = spreadsheet.getSheetByName('Last10');
  var query = 'SELECT ms.user.cell.value, ms.tag.cell.value, ms.sentiment.cell.value '
              +'FROM test_bq.ms '
              +'ORDER BY ms.data.cell.timestamp DESC '
              +'LIMIT 80'
  var row = 2;
  last10key = runQuerySQL(query)
  for (var i = 0; i < last10key.length; i=i+3) {
    user = last10key[i];
    tag = last10key[i+1];
    sentiment = last10key[i+2];
    sheet.getRange(row, 1).setValue(tag)
    sheet.getRange(row, 2).setValue(user)
    sheet.getRange(row, 3).setValue(sentiment)
    row++
  }
}

function runQuery(req){
  var ret = [];
  var projectId = 'slavayssiere-sandbox';

  var request = {
    query: req,
    location: 'EU',
    useLegacySql: 'true'
  };
                                
  var queryResults = BigQuery.Jobs.query(request, projectId);
  var jobId = queryResults.jobReference.jobId;

  // Check on status of the Query Job.
  var sleepTimeMs = 500;
  while (!queryResults.jobComplete) {
    Utilities.sleep(sleepTimeMs);
    sleepTimeMs *= 2;
    queryResults = BigQuery.Jobs.getQueryResults(projectId, jobId);
  }

  // Get all the rows of results.
  var rows = queryResults.rows;
  
  while (queryResults.pageToken) {
    queryResults = BigQuery.Jobs.getQueryResults(projectId, jobId, {
      pageToken: queryResults.pageToken
    });
    rows = rows.concat(queryResults.rows);
  }

  if (rows) {
    var num = 0;
    // Append the results.
    var data = new Array(rows.length);
    for (var i = 0; i < rows.length; i++) {
      var cols = rows[i].f;
      data[i] = new Array(cols.length);
      for (var j = 0; j < cols.length; j++) {
        ret[num] = cols[j].v;
        num++;
      }
    }

  } else {
    Logger.log('No rows returned.');
  }
  
  return ret;
}


function runQuerySQL(req){
  var ret = [];
  var projectId = 'slavayssiere-sandbox';

  var request = {
    query: req,
    location: 'EU',
    useLegacySql: 'true'
  };
                                
  var queryResults = BigQuery.Jobs.query(request, projectId);
  var jobId = queryResults.jobReference.jobId;

  // Check on status of the Query Job.
  var sleepTimeMs = 500;
  while (!queryResults.jobComplete) {
    Utilities.sleep(sleepTimeMs);
    sleepTimeMs *= 2;
    queryResults = BigQuery.Jobs.getQueryResults(projectId, jobId);
  }

  // Get all the rows of results.
  var rows = queryResults.rows;
  
  if (rows) {
    var num = 0;
    // Append the results.
    var data = new Array(rows.length);
    for (var i = 0; i < rows.length; i++) {
      var cols = rows[i].f;
      data[i] = new Array(cols.length);
      for (var j = 0; j < cols.length; j++) {
        ret[num] = cols[j].v;
        num++;
      }
    }

  } else {
    Logger.log('No rows returned.');
  }
  
  return ret;
}