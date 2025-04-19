package request

type JoinParty struct {
	// The join secret of the party.
	// 	example: "RzTMeBHZu3VoNEUpNQFSnMYpSNQgWQ2rYM4u3RHSPQIacMxE4KH63OYQEDAD0P0bnjZBYUFHB0I5wOx2xEiXOU9SKBxRZ5YcvYjjiZWhsGGss3vnRajvFn4trJgTnvz2"
	Secret string `json:"secret" binding:"required"`
}
