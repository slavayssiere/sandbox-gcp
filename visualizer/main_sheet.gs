function onOpen(e) {
  // Add a custom menu to the spreadsheet.
  SpreadsheetApp.getUi() // Or DocumentApp, SlidesApp, or FormApp.
      .createMenu('Get data')
      .addItem('Get users from BigTable', 'runQuery')
      .addSeparator()
      .addItem('Get data from other...', 'runQuery')
      .addToUi();
}

/**
 * Runs a BigQuery query and logs the results in a spreadsheet.
 * Helper: https://developers.google.com/apps-script/guides/services/advanced
 */
function runQuery() {
  var projectId = 'slavayssiere-sandbox';

  var request = {
    query: 'SELECT ms.user.cell.value FROM test_bq.ms LIMIT 1000',
    location: 'EU'
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
    var spreadsheet = SpreadsheetApp.getActiveSpreadsheet()
    var sheet = spreadsheet.getActiveSheet()

    // Append the headers.
    var headers = queryResults.schema.fields.map(function(field) {
      return field.name;
    });
    sheet.appendRow(headers);

    // Append the results.
    var data = new Array(rows.length);
    for (var i = 0; i < rows.length; i++) {
      var cols = rows[i].f;
      data[i] = new Array(cols.length);
      for (var j = 0; j < cols.length; j++) {
        data[i][j] = cols[j].v;
      }
    }
    sheet.getRange(2, 1, rows.length, headers.length).setValues(data);

  } else {
    Logger.log('No rows returned.');
  }
}
