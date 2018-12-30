const pusher = new Pusher('ba844c624003f02c6c0f', {
  cluster: 'ap1',
  encrypted: true
});

const channel = pusher.subscribe('tootCountHourly');

channel.bind('addNumber', data => {
if (newLineChart.data.labels.length > 15) {
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
      label: "草莓馅小时嘟嘟量",
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
      label: "草莓馅日嘟嘟量",
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