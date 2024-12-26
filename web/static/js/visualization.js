class ChordVisualizer {
    constructor() {
        this.width = 800;
        this.height = 600;
        this.radius = Math.min(this.width, this.height) / 2 - 50;
        
        this.svg = d3.select("#chord-ring")
            .append("svg")
            .attr("width", this.width)
            .attr("height", this.height)
            .append("g")
            .attr("transform", `translate(${this.width/2},${this.height/2})`);
            
        this.nodes = [];
        this.initializeWebSocket();
    }
    
    initializeWebSocket() {
        const ws = new WebSocket(`ws://${window.location.host}/ws`);
        
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.updateVisualization(data);
        };
        
        ws.onerror = (error) => {
            console.error("WebSocket error:", error);
        };
        
        ws.onclose = () => {
            console.log("WebSocket connection closed");
            setTimeout(() => this.initializeWebSocket(), 5000);
        };
    }
    
    updateVisualization(data) {
        this.nodes = data;
        

        this.svg.selectAll("*").remove();
        this.svg.append("circle")
            .attr("r", this.radius)
            .attr("fill", "none")
            .attr("stroke", "#e2e8f0")
            .attr("stroke-width", 1);
        const nodePositions = this.calculateNodePositions();
        this.drawConnections(nodePositions);
        this.drawFingerConnections(nodePositions);
        this.drawNodes(nodePositions);
    }
    
    calculateNodePositions() {
        const positions = new Map();
        const angleStep = (2 * Math.PI) / Math.pow(2, 160); // For 160-bit ID space
        
        this.nodes.forEach(node => {
            const id = BigInt("0x" + node.id);
            const angle = Number(id) * angleStep;
            positions.set(node.id, {
                x: this.radius * Math.cos(angle),
                y: this.radius * Math.sin(angle),
                angle: angle
            });
        });
        
        return positions;
    }
    
    drawConnections(positions) {
        this.nodes.forEach(node => {
            if (node.successor) {
                const source = positions.get(node.id);
                const target = positions.get(node.successor);
                
                if (source && target) {
                    this.svg.append("path")
                        .attr("class", "link")
                        .attr("d", this.generatePath(source, target))
                        .attr("fill", "none");
                }
            }
        });
    }
    
    drawFingerConnections(positions) {
        this.nodes.forEach(node => {
            node.fingers.forEach(fingerId => {
                const source = positions.get(node.id);
                const target = positions.get(fingerId);
                
                if (source && target) {
                    this.svg.append("path")
                        .attr("class", "finger-link")
                        .attr("d", this.generatePath(source, target))
                        .attr("fill", "none");
                }
            });
        });
    }
    
    drawNodes(positions) {
        this.nodes.forEach(node => {
            const pos = positions.get(node.id);
            if (pos) {
                const nodeGroup = this.svg.append("g")
                    .attr("class", "node")
                    .attr("transform", `translate(${pos.x},${pos.y})`);
                
                nodeGroup.append("circle")
                    .attr("r", 6);
                
                nodeGroup.append("text")
                    .attr("x", 10)
                    .attr("y", 4)
                    .text(node.address);
            }
        });
    }
    
    generatePath(source, target) {
        const dx = target.x - source.x;
        const dy = target.y - source.y;
        const dr = Math.sqrt(dx * dx + dy * dy);
        
        return `M${source.x},${source.y}A${dr},${dr} 0 0,1 ${target.x},${target.y}`;
    }
}

window.addEventListener('load', () => {
    new ChordVisualizer();
});