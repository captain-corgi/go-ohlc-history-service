# Test result

Experimental

## Config

### Docker machine

    MacBook Air (13-inch, M1, 2020)
    Chip Apple M1
    Memory 8 GB

### DB (MySQL)

    Windows 11
    Memory 16 GB
    Chip Intel Core i5-9300H
    Version: 5.7.28
    Storage: 50 GB

## Result

| No  | size           | rows              | time         | result    |
|-----|----------------|-------------------|--------------|-----------|
| 1   | Less than 10MB | less than 10000   | very fast    | real      |
| 2   | 400 MB         | 5,242,880         | 170,774 ms   | real      |
| 3   | 1.72 GB        | 20,971,520        | 676,052 ms   | real      |
| 4   | 2.5 GB         | 31,931,040        | 17,3146 min  | estimated |
| 5   | 3.5 GB         | 43,931,040        | 30.648 min   | estimated |
| 6   | 100GB          | 1,038,876,000     | 63.98 min    | estimated |
| 7   | 200GB          | 2,038,876,000     | 130.648 min  | estimated |
| 8   | 300GB          | 3,038,876,000     | 230.648 min  | estimated |
| 9   | 500GB          | 5,038,876,000     | 397.3146 min | estimated |
| 10  | 1 TB           | 10,038,876,000    | 13.2886 hour | estimated |
| 11  | 2 TB           | 20,038,876,000    | 1.1325 day   | estimated |
| 12  | 3 TB           | 30,038,876,000    | 2.29167 day  | estimated |
| 13  | 4 TB           | 40,038,876,000    | 3.5 day      | estimated |
