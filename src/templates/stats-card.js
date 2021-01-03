const Card = require("./card");
const d3 = require("d3");
const moment = require("moment");

function createStatsCard(title, statsData, theme) {
  let card = new Card(title, 230, 200, theme);
  let svg = card.getSVG();

  //draw icon
  const panel = svg.append("g").attr("transform", `translate(30,20)`);
  let labelHeight = 14;
  panel
    .selectAll(null)
    .data(statsData)
    .enter()
    .append("g")
    .attr("transform", (d) => {
      let y = labelHeight * d.index * 1.8;
      return `translate(0,${y})`;
    })
    .attr("width", labelHeight)
    .attr("height", labelHeight)
    .attr("fill", theme.icon_color)
    .html((d) => d.icon);

  //draw text
  panel
    .selectAll(null)
    .data(statsData)
    .enter()
    .append("text")
    .text((d) => {
      return `${d.name}`;
    })
    .attr("x", labelHeight * 1.5)
    .attr("y", (d) => labelHeight * d.index * 1.8 + labelHeight)
    .style("fill", theme.text_color)
    .style("font-size", `${labelHeight}px`);

  panel
    .selectAll(null)
    .data(statsData)
    .enter()
    .append("text")
    .text((d) => {
      return `${d.value}`;
    })
    .attr("x", 140)
    .attr("y", (d) => labelHeight * d.index * 1.8 + labelHeight)
    .style("fill", theme.text_color)
    .style("font-size", `${labelHeight}px`);

  return card.toString();
}

module.exports = createStatsCard;
