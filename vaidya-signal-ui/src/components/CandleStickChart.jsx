import Highcharts from "highcharts/highstock";
import HighchartsReact from "highcharts-react-official";

export default function CandlestickChart({ ticker }) {
  const options = {
    title: {
      text: ticker,
    },
    series: [{
      type: "candlestick",
      data: [
        [Date.UTC(2022, 5, 2), 5, 9, 4, 7],
        [Date.UTC(2022, 5, 3), 7, 16, 7, 15],
        [Date.UTC(2022, 5, 4), 9, 12, 3, 5],
      ],
    }, {
      type: "line",
      name: "Close",
      data: [[Date.UTC(2022, 5, 2), 7], [Date.UTC(2022, 5, 3), 15], [Date.UTC(2022, 5, 4), 5]],
    }],
    xAxis: {
      type: "datetime",
      plotLines: [{
        color: "red", // Color of the line
        width: 2, // Width of the line
        value: Date.UTC(2022, 5, 3), // X-axis value where the line should be placed
        zIndex: 5, // Z-index of the line, to make sure it appears above the candlestick series
      }],
    },
  };

  return (
    <div className="w-full">
      <HighchartsReact
        highcharts={Highcharts}
        options={options}
      />
    </div>
  );
}
