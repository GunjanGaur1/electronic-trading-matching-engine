import { useEffect, useMemo, useState } from "react";
import "./App.css";

const initialAsks = [
  { price: 102.0, quantity: 80, orders: 2 },
  { price: 101.5, quantity: 45, orders: 1 },
  { price: 101.0, quantity: 50, orders: 1 },
];

const initialBids = [
  { price: 100.5, quantity: 35, orders: 2 },
  { price: 100.0, quantity: 50, orders: 2 },
  { price: 99.5, quantity: 100, orders: 1 },
];

const initialTrades = [
  { id: 1, time: "00:42:18.221", side: "BUY", quantity: 40, price: 102.0 },
  { id: 2, time: "00:42:16.904", side: "SELL", quantity: 20, price: 101.5 },
  { id: 3, time: "00:42:14.731", side: "BUY", quantity: 30, price: 101.0 },
];

function App() {
  const [asks, setAsks] = useState(initialAsks);
  const [bids, setBids] = useState(initialBids);
  const [trades, setTrades] = useState(initialTrades);

  const [side, setSide] = useState("BUY");
  const [orderType, setOrderType] = useState("LIMIT");
  const [price, setPrice] = useState("101.00");
  const [quantity, setQuantity] = useState("100");

  const [orderId, setOrderId] = useState(1247);
  const [ordersProcessed, setOrdersProcessed] = useState(1284732);
  const [latency, setLatency] = useState(42);

  useEffect(() => {
    const interval = setInterval(() => {
      setLatency(Math.floor(35 + Math.random() * 20));
      setOrdersProcessed((value) => value + Math.floor(Math.random() * 4));
    }, 1200);

    return () => clearInterval(interval);
  }, []);

  const bestAsk = asks[0]?.price ?? 0;
  const bestBid = bids[0]?.price ?? 0;

  const spread = useMemo(() => {
    if (!bestAsk || !bestBid) return "0.00";
    return (bestAsk - bestBid).toFixed(2);
  }, [bestAsk, bestBid]);

  const midPrice = useMemo(() => {
    if (!bestAsk || !bestBid) return "0.00";
    return ((bestAsk + bestBid) / 2).toFixed(2);
  }, [bestAsk, bestBid]);

  function submitOrder(event) {
    event.preventDefault();

    const qty = Number(quantity);
    const px = Number(price);

    if (!qty || qty <= 0) return;

    const newOrderId = orderId + 1;
    setOrderId(newOrderId);

    const now = new Date();

    const timestamp = now
      .toLocaleTimeString("en-US", {
        hour12: false,
      })
      .concat(".", String(now.getMilliseconds()).padStart(3, "0"));

    /*
      Demo-only local matching behavior.

      Later this will be replaced by:
        React -> Go API/WebSocket -> Matching Engine
    */

    if (orderType === "MARKET") {
      const newTrade = {
        id: Date.now(),
        time: timestamp,
        side,
        quantity: qty,
        price: side === "BUY" ? bestAsk : bestBid,
      };

      setTrades((current) => [newTrade, ...current].slice(0, 8));
    } else {
      const newLevel = {
        price: px,
        quantity: qty,
        orders: 1,
      };

      if (side === "BUY") {
        setBids((current) =>
          [...current, newLevel]
            .sort((a, b) => b.price - a.price)
            .slice(0, 6)
        );
      } else {
        setAsks((current) =>
          [...current, newLevel]
            .sort((a, b) => a.price - b.price)
            .slice(0, 6)
        );
      }
    }

    setQuantity("");
  }

  return (
    <div className="app">
      <header className="topbar">
        <div className="brand">
          <div className="brand-mark">E</div>

          <div>
            <div className="brand-name">ENGINE</div>
            <div className="brand-subtitle">
              ELECTRONIC TRADING SYSTEM
            </div>
          </div>
        </div>

        <div className="market-info">
          <div>
            <span className="muted">MARKET</span>
            <strong>BTC / USD</strong>
          </div>

          <div>
            <span className="muted">MID</span>
            <strong>${midPrice}</strong>
          </div>

          <div>
            <span className="muted">SPREAD</span>
            <strong>${spread}</strong>
          </div>

          <div className="status">
            <span className="status-dot"></span>
            ENGINE ONLINE
          </div>
        </div>
      </header>

      <main className="dashboard">
        <section className="stats">
          <Stat
            label="ORDERS PROCESSED"
            value={ordersProcessed.toLocaleString()}
            change="+2.4%"
          />

          <Stat
            label="MATCHING LATENCY"
            value={`${latency} μs`}
            change="LIVE"
          />

          <Stat
            label="BEST BID"
            value={`$${bestBid.toFixed(2)}`}
            change="BUY"
          />

          <Stat
            label="BEST ASK"
            value={`$${bestAsk.toFixed(2)}`}
            change="SELL"
          />
        </section>

        <section className="main-grid">
          <div className="panel orderbook-panel">
            <PanelHeader
              title="ORDER BOOK"
              subtitle="PRICE-TIME PRIORITY"
            />

            <div className="book-header">
              <span>PRICE</span>
              <span>SIZE</span>
              <span>ORDERS</span>
            </div>

            <div className="book asks">
              {asks.map((level, index) => (
                <BookRow
                  key={`${level.price}-${index}`}
                  level={level}
                  type="ask"
                />
              ))}
            </div>

            <div className="spread-row">
              <span>SPREAD</span>
              <strong>${spread}</strong>
            </div>

            <div className="book bids">
              {bids.map((level, index) => (
                <BookRow
                  key={`${level.price}-${index}`}
                  level={level}
                  type="bid"
                />
              ))}
            </div>
          </div>

          <div className="panel order-panel">
            <PanelHeader
              title="PLACE ORDER"
              subtitle={`ORDER #${orderId + 1}`}
            />

            <div className="toggle-group">
              <button
                className={side === "BUY" ? "active buy" : ""}
                onClick={() => setSide("BUY")}
              >
                BUY
              </button>

              <button
                className={side === "SELL" ? "active sell" : ""}
                onClick={() => setSide("SELL")}
              >
                SELL
              </button>
            </div>

            <div className="type-group">
              <button
                className={orderType === "LIMIT" ? "selected" : ""}
                onClick={() => setOrderType("LIMIT")}
              >
                LIMIT
              </button>

              <button
                className={orderType === "MARKET" ? "selected" : ""}
                onClick={() => setOrderType("MARKET")}
              >
                MARKET
              </button>
            </div>

            <form onSubmit={submitOrder}>
              {orderType === "LIMIT" && (
                <label>
                  <span>PRICE</span>

                  <div className="input-wrap">
                    <span>$</span>

                    <input
                      value={price}
                      onChange={(e) => setPrice(e.target.value)}
                      type="number"
                      step="0.01"
                    />
                  </div>
                </label>
              )}

              <label>
                <span>QUANTITY</span>

                <div className="input-wrap">
                  <input
                    value={quantity}
                    onChange={(e) => setQuantity(e.target.value)}
                    type="number"
                    placeholder="100"
                  />

                  <span>SHARES</span>
                </div>
              </label>

              <div className="order-summary">
                <div>
                  <span>ORDER TYPE</span>
                  <strong>{orderType}</strong>
                </div>

                <div>
                  <span>TIME IN FORCE</span>
                  <strong>DAY</strong>
                </div>

                <div>
                  <span>ROUTING</span>
                  <strong>DIRECT</strong>
                </div>
              </div>

              <button className={`submit ${side.toLowerCase()}`}>
                SUBMIT {side} ORDER
              </button>
            </form>
          </div>

          <div className="panel trades-panel">
            <PanelHeader
              title="RECENT TRADES"
              subtitle="EXECUTION FEED"
            />

            <div className="trade-header">
              <span>TIME</span>
              <span>SIDE</span>
              <span>SIZE</span>
              <span>PRICE</span>
            </div>

            <div className="trades">
              {trades.map((trade) => (
                <div className="trade-row" key={trade.id}>
                  <span className="time">{trade.time}</span>

                  <span
                    className={
                      trade.side === "BUY"
                        ? "trade-buy"
                        : "trade-sell"
                    }
                  >
                    {trade.side}
                  </span>

                  <span>{trade.quantity}</span>

                  <strong>${trade.price.toFixed(2)}</strong>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section className="bottom-grid">
          <div className="panel architecture-panel">
            <PanelHeader
              title="ENGINE STATUS"
              subtitle="SYSTEM ARCHITECTURE"
            />

            <div className="architecture">
              <ArchitectureNode
                title="ORDER INTAKE"
                value="1.28M"
                status="ACTIVE"
              />

              <div className="arrow">→</div>

              <ArchitectureNode
                title="ORDER BOOK"
                value={`${bids.length + asks.length} LEVELS`}
                status="ACTIVE"
              />

              <div className="arrow">→</div>

              <ArchitectureNode
                title="MATCH ENGINE"
                value={`${latency} μs`}
                status="ACTIVE"
              />

              <div className="arrow">→</div>

              <ArchitectureNode
                title="TRADE OUTPUT"
                value="REALTIME"
                status="ACTIVE"
              />
            </div>
          </div>

          <div className="panel tech-panel">
            <PanelHeader
              title="SYSTEM"
              subtitle="RUNTIME"
            />

            <div className="tech-list">
              <Tech name="MATCHING ENGINE" value="GO" />
              <Tech name="CONCURRENCY" value="GOROUTINES" />
              <Tech name="SYNCHRONIZATION" value="RWMutex" />
              <Tech name="CORE" value="C++17" />
              <Tech name="PROTOCOL" value="TCP/IP" />
            </div>
          </div>
        </section>
      </main>

      <footer>
        <span>ENGINE v1.0.0</span>
        <span>LOW-LATENCY MATCHING SYSTEM</span>
        <span>STATUS: OPERATIONAL</span>
      </footer>
    </div>
  );
}

function Stat({ label, value, change }) {
  return (
    <div className="stat">
      <span>{label}</span>

      <div>
        <strong>{value}</strong>
        <small>{change}</small>
      </div>
    </div>
  );
}

function PanelHeader({ title, subtitle }) {
  return (
    <div className="panel-header">
      <div>
        <h2>{title}</h2>
        <span>{subtitle}</span>
      </div>

      <div className="panel-indicator"></div>
    </div>
  );
}

function BookRow({ level, type }) {
  const width = Math.min(100, (level.quantity / 100) * 100);

  return (
    <div className={`book-row ${type}`}>
      <div
        className="depth"
        style={{ width: `${width}%` }}
      ></div>

      <span className="price">
        ${level.price.toFixed(2)}
      </span>

      <span>{level.quantity}</span>

      <span>{level.orders}</span>
    </div>
  );
}

function ArchitectureNode({ title, value, status }) {
  return (
    <div className="architecture-node">
      <span>{title}</span>
      <strong>{value}</strong>

      <small>
        <i></i>
        {status}
      </small>
    </div>
  );
}

function Tech({ name, value }) {
  return (
    <div className="tech-row">
      <span>{name}</span>
      <strong>{value}</strong>
    </div>
  );
}

export default App;