# qwakes

Software for [Q-Hawkes process](https://arxiv.org/abs/1509.07710) simulation and inference on stock market data. The package is written primarily in Go and uses the [stochadex](https://github.com/umbralcalc/stochadex) as the simulation engine which powers all of the essential tools developed in this repo.

**Plan**

- Use the [Polygon Go client](https://github.com/polygon-io/client-go) to create functions for reading stock market data and loading it into `simulator.StateTimeStorage` objects 
- Create a custom Q-Hawkes simulation iterator and some convenience functions for running them with the market data
- Develop a custom online inference methodology and mechanism for cross-validation against the data, then create some custom functions for this
- Create a portfolio agent iterator which acquires positions in the market and affects the market liquidity (and possibly price) through its actions via the Hawkes model
- Develop some more convenience functions for backtesting portfolio configurations 
  - Obviously need to reconcile simulation disparity with historical data...
  - So is it possible to use the historical trade data as possible actions/shifts in portfolio positions?
