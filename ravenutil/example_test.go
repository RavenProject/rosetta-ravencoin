package ravenutil_test

import (
	"fmt"
	"math"

	"github.com/RavenProject/rosetta-ravencoin/ravenutil"
)

func ExampleAmount() {

	a := ravenutil.Amount(0)
	fmt.Println("Zero Satoshi:", a)

	a = ravenutil.Amount(1e8)
	fmt.Println("100,000,000 Satoshis:", a)

	a = ravenutil.Amount(1e5)
	fmt.Println("100,000 Satoshis:", a)
	// Output:
	// Zero Satoshi: 0 RVN
	// 100,000,000 Satoshis: 1 RVN
	// 100,000 Satoshis: 0.001 RVN
}

func ExampleNewAmount() {
	amountOne, err := ravenutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := ravenutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := ravenutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := ravenutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 RVN
	// 0.01234567 RVN
	// 0 RVN
	// invalid bitcoin amount
}

func ExampleAmount_unitConversions() {
	amount := ravenutil.Amount(44433322211100)

	fmt.Println("Satoshi to kRVN:", amount.Format(ravenutil.AmountKiloBTC))
	fmt.Println("Satoshi to RVN:", amount)
	fmt.Println("Satoshi to MilliRVN:", amount.Format(ravenutil.AmountMilliBTC))
	fmt.Println("Satoshi to MicroRVN:", amount.Format(ravenutil.AmountMicroBTC))
	fmt.Println("Satoshi to Satoshi:", amount.Format(ravenutil.AmountSatoshi))

	// Output:
	// Satoshi to kRVN: 444.333222111 kRVN
	// Satoshi to RVN: 444333.222111 RVN
	// Satoshi to MilliRVN: 444333222.111 mRVN
	// Satoshi to MicroRVN: 444333222111 μRVN
	// Satoshi to Satoshi: 44433322211100 Satoshi
}
