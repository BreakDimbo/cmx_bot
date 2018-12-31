const pusher = new Pusher('ba844c624003f02c6c0f', {
  cluster: 'ap1',
  encrypted: true
});

const channel = pusher.subscribe('tootCountHourly');

var HttpClient = function() {
  this.get = function(aUrl, aCallback) {
      var anHttpRequest = new XMLHttpRequest();
      anHttpRequest.onreadystatechange = function() { 
          if (anHttpRequest.readyState == 4 && anHttpRequest.status == 200)
              aCallback(anHttpRequest.responseText);
      }

      anHttpRequest.open( "GET", aUrl, true );            
      anHttpRequest.send( null );
  }
}

var client = new HttpClient();
client.get('http://67.216.197.45:8085/tootCountHourly', function(response) {
  var res = JSON.parse(response)
  res = res.filter(function(e) {return e != null})
  res = res.reverse()
  res.forEach(element => {
    newLineChart.data.labels.push(element.Time);
    newLineChart.data.datasets[0].data.push(element.Count);
    newLineChart.update();
  });
});

channel.bind('addNumber', data => {
if (newLineChart.data.labels.length > 24) {
  newLineChart.data.labels.shift();  
  newLineChart.data.datasets[0].data.shift();
}

newLineChart.data.labels.push(data.Time);
newLineChart.data.datasets[0].data.push(data.Count);
newLineChart.update();
});

function renderChart(userVisitsData) {
  var ctx = document.getElementById("hourcount").getContext("2d");

  var options = {};

  newLineChart = new Chart(ctx, {
    type: "line",
    data: userVisitsData,
    options: options
  });
}

var chartConfig = {
labels: [],
datasets: [
   {
      label: "草莓县县民小时嘟嘟量",
      fill: false,
      lineTension: 0.1,
      backgroundColor: "rgba(75,192,192,0.4)",
      borderColor: "rgba(75,192,192,1)",
      borderCapStyle: 'butt',
      borderDash: [],
      borderDashOffset: 0.0,
      borderJoinStyle: 'miter',
      pointBorderColor: "rgba(75,192,192,1)",
      pointBackgroundColor: "#fff",
      pointBorderWidth: 1,
      pointHoverRadius: 5,
      pointHoverBackgroundColor: "rgba(75,192,192,1)",
      pointHoverBorderColor: "rgba(220,220,220,1)",
      pointHoverBorderWidth: 2,
      pointRadius: 1,
      pointHitRadius: 10,
      data: [],
      spanGaps: false,
   }
]
};

const channelDaily = pusher.subscribe('tootCountDaily');

client.get('http://67.216.197.45:8085/tootCountDaily', function(response) {
  var res = JSON.parse(response)
  res = res.filter(function(e) {return e != null})
  res = res.reverse()
  res.forEach(element => {
    newDailyLineChart.data.labels.push(element.Time);
    newDailyLineChart.data.datasets[0].data.push(element.Count);
    newDailyLineChart.update();
  });
});

channelDaily.bind('addNumber', data => {
if (newDailyLineChart.data.labels.length > 15) {
  newDailyLineChart.data.labels.shift();  
  newDailyLineChart.data.datasets[0].data.shift();
}

newDailyLineChart.data.labels.push(data.Time);
newDailyLineChart.data.datasets[0].data.push(data.Count);
newDailyLineChart.update();
});

function renderDailyChart(userVisitsData) {
  var ctx = document.getElementById("dailycount").getContext("2d");

  var options = {};

  newDailyLineChart = new Chart(ctx, {
    type: "line",
    data: userVisitsData,
    options: options
  });
}

var chartConfigDaily = {
labels: [],
datasets: [
   {
      label: "草莓县县民日嘟嘟量",
      fill: false,
      lineTension: 0.1,
      backgroundColor: "rgba(75,192,192,0.4)",
      borderColor: "rgba(75,192,192,1)",
      borderCapStyle: 'butt',
      borderDash: [],
      borderDashOffset: 0.0,
      borderJoinStyle: 'miter',
      pointBorderColor: "rgba(75,192,192,1)",
      pointBackgroundColor: "#fff",
      pointBorderWidth: 1,
      pointHoverRadius: 5,
      pointHoverBackgroundColor: "rgba(75,192,192,1)",
      pointHoverBorderColor: "rgba(220,220,220,1)",
      pointHoverBorderWidth: 2,
      pointRadius: 1,
      pointHitRadius: 10,
      data: [],
      spanGaps: false,
   }
]
};

renderDailyChart(chartConfigDaily)
renderChart(chartConfig)