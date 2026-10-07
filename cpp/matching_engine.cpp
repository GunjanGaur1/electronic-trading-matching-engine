#include <iostream>
#include <map>
#include <queue>
#include <string>

enum class Side {
    BUY,
    SELL
};

enum class OrderType {
    LIMIT,
    MARKET
};

struct Order {
    int id;
    Side side;
    OrderType type;
    double price;
    int quantity;
};

struct Trade {
    int buyOrderId;
    int sellOrderId;
    double price;
    int quantity;
};

class MatchingEngine {

private:

    std::map<double, std::queue<Order>, std::greater<double>> bids;
    std::map<double, std::queue<Order>> asks;

public:

    void addOrder(Order order) {

        if (order.side == Side::BUY) {
            matchBuy(order);
        } else {
            matchSell(order);
        }
    }

private:

    void matchBuy(Order& incoming) {

        while (
            incoming.quantity > 0 &&
            !asks.empty()
        ) {

            auto bestAsk = asks.begin();

            if (
                incoming.type == OrderType::LIMIT &&
                incoming.price < bestAsk->first
            ) {
                break;
            }

            Order& resting = bestAsk->second.front();

            int quantity =
                std::min(
                    incoming.quantity,
                    resting.quantity
                );

            std::cout
                << "TRADE "
                << quantity
                << " @ "
                << resting.price
                << std::endl;

            incoming.quantity -= quantity;
            resting.quantity -= quantity;

            if (resting.quantity == 0) {
                bestAsk->second.pop();

                if (bestAsk->second.empty()) {
                    asks.erase(bestAsk);
                }
            }
        }

        if (
            incoming.quantity > 0 &&
            incoming.type == OrderType::LIMIT
        ) {
            bids[incoming.price].push(incoming);
        }
    }

    void matchSell(Order& incoming) {

        while (
            incoming.quantity > 0 &&
            !bids.empty()
        ) {

            auto bestBid = bids.begin();

            if (
                incoming.type == OrderType::LIMIT &&
                incoming.price > bestBid->first
            ) {
                break;
            }

            Order& resting = bestBid->second.front();

            int quantity =
                std::min(
                    incoming.quantity,
                    resting.quantity
                );

            std::cout
                << "TRADE "
                << quantity
                << " @ "
                << resting.price
                << std::endl;

            incoming.quantity -= quantity;
            resting.quantity -= quantity;

            if (resting.quantity == 0) {
                bestBid->second.pop();

                if (bestBid->second.empty()) {
                    bids.erase(bestBid);
                }
            }
        }

        if (
            incoming.quantity > 0 &&
            incoming.type == OrderType::LIMIT
        ) {
            asks[incoming.price].push(incoming);
        }
    }
};

int main() {

    MatchingEngine engine;

    engine.addOrder({
        1,
        Side::SELL,
        OrderType::LIMIT,
        100.0,
        50
    });

    engine.addOrder({
        2,
        Side::BUY,
        OrderType::LIMIT,
        100.0,
        30
    });

    return 0;
}