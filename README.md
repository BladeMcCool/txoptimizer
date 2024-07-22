# Transaction Optimizer (Knapsack Problem)

### FingerPrint take home assessment - 2024-05-19 (May 19 2024)

[Full problem description can be found here](challenge_source/fp%20transaction%20processing%20task.md)  
TLDR: "Find the best set of transactions we can send within the allotted time"  
I spent roughly a weekend researching this, understanding the logic and settling on an implementation to submit.

## My Approach

My initial thoughts on this challenge were that it sounded extremely straightforward and that it could easily be solved by sorting things by the highest value per ms of tx time and then take the transactions which would fit in the time window off the top of that sorted list until there was no more time left. 

However, realizing that this was probably too naive I decided to do some research. Upon researching, I discovered that this is the knapsack problem, which has some documented solutions including algorithms that appear to produce the correct output for our problem here.

I tried out several implementations before settling on the single dimensional data collection one. This one is optimal from what I looked at because it gets the correct results compared to my naive first attempt, and uses a bit less memory than the multidimensional value tracker that I also tried out. There are probably more optimal ways to solve it, as it seems this is a well studied problem, but I spent most of my time trying to better understand things. We go through all the items, and all the timeslots, and track if we will get a better value by including the item or not, in a array of timeslots basically. We track if we think we would keep that item at that point in time and then when we're done we go through the keep list jumping back down the timeslots to find the tx we actually will select.

## A summary of the knapsack problem is documented by youtuber Mindez here:

https://www.youtube.com/watch?v=EH6h7WA7sDw

This is similar to one of the implementations I tried.

## The answers 

For max time 50 ms:
Max USD value: 4139.43

For max time 60 ms:
Max USD value: 4675.71

For max time 90 ms:
Max USD value: 6972.29

For max time 1000 ms:
Max USD value: 35471.81
