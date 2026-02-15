import { useEffect, useState } from "react";
import {
  Chart as ChartJS,
  LineElement,
  PointElement,
  LinearScale,
  TimeScale,
  CategoryScale,
  Tooltip,
  Legend,
} from "chart.js";
import { Line } from "react-chartjs-2";
import "chartjs-adapter-date-fns";
import "./App.css";

ChartJS.register(
  LineElement,
  PointElement,
  LinearScale,
  TimeScale,
  CategoryScale,
  Tooltip,
  Legend,
);

function App() {
  const [chartData, setChartData] = useState(null);

  useEffect(() => {
    fetch("/data.json")
      .then((res) => res.json())
      .then((data) => {
        const vibrantColors = [
          "#FF6B6B",
          "#4ECDC4",
          "#FFD93D",
          "#6C5CE7",
          "#00B894",
          "#E84393",
        ];

        const datasets = Object.keys(data).map((user, index) => ({
          label: user,
          data: data[user].map((entry) => ({
            x: entry.timestamp,
            y: entry.total,
          })),
          borderColor: vibrantColors[index % vibrantColors.length],
          backgroundColor: vibrantColors[index % vibrantColors.length] + "33",
          borderWidth: 3,
          tension: 0.35,
          pointRadius: 4,
          pointHoverRadius: 6,
          fill: true,
        }));

        setChartData({ datasets });
      });
  }, []);

  if (!chartData) return <div style={{ padding: "40px" }}>Loading...</div>;

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        labels: {
          color: "#c9d1d9",
          padding: 20,
          boxWidth: 30,
        },
      },
    },
    scales: {
      x: {
        type: "time",
        title: { display: true, text: "Time", color: "#c9d1d9" },
        time: { unit: "day" },
        ticks: { color: "#8b949e" },
        grid: { color: "rgba(255,255,255,0.06)" },
      },
      y: {
        title: { display: true, text: "Total Solved", color: "#c9d1d9" },
        ticks: { color: "#8b949e", precision: 0 },
        grid: { color: "rgba(255,255,255,0.06)" },
      },
    },
  };

  return (
    <div className="page">
      <div className="title">LeetCode Progress</div>
      <div className="card">
        <div className="chart-container">
          <Line data={chartData} options={options} />
        </div>
      </div>
    </div>
  );
}

export default App;
