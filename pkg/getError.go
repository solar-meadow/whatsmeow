package pkg

import "log"

func getErrorMessage(errCode int) string {
	var str string

	switch errCode {
	case 0:
		return str
	case 1:
		str = ErrInvalidNumber1
	case 6:
		str = ErrSubscriberNotInNet
	case 11:
		str = ErrServiceNotConnected
	case 12:
		str = ErrInvalidPhoneNumber
	case 13:
		str = ErrSubscriberBlocked
	case 21:
		str = ErrNoServiceSupport
	case 200:
		str = ErrVirtualSending
	case 219:
		str = ErrSimCardReplacement
	case 220:
		str = ErrOperatorQueueOverflow
	case 237:
		str = ErrSubscriberNotAnswering
	case 238:
		str = ErrNoTemplate
	case 239:
		str = ErrForbiddenIPAddress
	case 240:
		str = ErrSubscriberBusy
	case 241:
		str = ErrConversionError
	case 242:
		str = ErrAnsweringMachineDetected
	case 243:
		str = ErrUnregisteredSenderID
	case 244:
		str = ErrRejectedByOperator
	case 245:
		str = ErrInvalidFormatNumber
	case 246:
		str = ErrNumberNotAllowedBySettings
	case 247:
		str = ErrDailyMessageLimitExceeded
	case 248:
		str = ErrNoRoute
	case 249:
		str = ErrInvalidFormatNumber249
	case 250:
		str = ErrNumberProhibitedBySettings
	case 251:
		str = ErrExceedDailyLimitPerNumber
	case 252:
		str = ErrNumberProhibited
	case 253:
		str = ErrSpamFilterForbidden
	case 254:
		str = ErrUnregisteredSenderID255
	case 255:
		str = ErrOperatorRejected
	default:
		str = "Неизвестная ошибка"
		log.Println(errCode)
	}
	return str
}
