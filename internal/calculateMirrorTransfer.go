package internal

func CalculateMirrorTransfer(walletIdSenderArray []string, walletIdRecipientArray []string, amount []int) ([]string, []string, []int) {

	intersections := findMirrorTransfer(walletIdSenderArray, walletIdRecipientArray)

	if len(intersections) > 0 {
		sendersCalcMirror, recipientsCalcMirror, newAmount := calculateMirrorTransactions(intersections, amount)
		return generateNewCalculateTransfer(walletIdSenderArray, walletIdRecipientArray, amount, intersections, sendersCalcMirror, recipientsCalcMirror, newAmount)
	} else {
		return walletIdSenderArray, walletIdRecipientArray, amount
	}

}

func findMirrorTransfer(walletIdSenderArray []string, walletIdRecipientArray []string) []int {
	added := make(map[int]bool)

	// Ищем пересечения и сохраняем их индексы в новый массив
	var intersections []int
	for i, elem1 := range walletIdSenderArray {
		for j, elem2 := range walletIdRecipientArray {
			if elem1 == elem2 && walletIdSenderArray[j] == walletIdRecipientArray[i] && !added[i] && !added[j] {
				intersections = append(intersections, i, j)
				added[i] = true
				added[j] = true
			}
		}
	}

	return intersections
}

// возвращает уникальные индексы c с массива senders зеркальных транзакция с конечной суммой трансфера, проведя их рассчет
func calculateMirrorTransactions(mirrorIndex []int, amount []int) ([]int, []int, []int) {

	var senders []int
	var recipients []int
	var newAmount []int

	for i := 0; i <= len(mirrorIndex)-2; i = i + 2 {
		if amount[i] < amount[i+1] {
			senders = append(senders, mirrorIndex[i])
			recipients = append(recipients, mirrorIndex[i+1])
			newAmount = append(newAmount, amount[i+1]-amount[i])

		} else {
			senders = append(senders, mirrorIndex[i+1])
			recipients = append(recipients, mirrorIndex[i])
			newAmount = append(newAmount, amount[i]-amount[i+1])
		}
	}
	return senders, recipients, newAmount
}
func generateNewCalculateTransfer(walletIdSenderArray []string, walletIdRecipientArray []string, amount []int, mirrorIndex []int, sendersNew []int, recipientsNew []int, newMirrorAmount []int) ([]string, []string, []int) {
	// Создаем новые срезы для хранения элементов, которые останутся после удаления
	var newSenderArray []string
	var newRecipientArray []string
	var newAmountArray []int

	for i, sender := range walletIdSenderArray {
		// Если индекс элемента не находится в массиве mirrorIndex,
		// то добавляем этот элемент в новый срез
		if !contains(mirrorIndex, i) {
			newSenderArray = append(newSenderArray, sender)
			newRecipientArray = append(newRecipientArray, walletIdRecipientArray[i])
			newAmountArray = append(newAmountArray, amount[i])
		}
	}
	//добавляем перерасчитанные зеркальные трансферы
	for i := range sendersNew {
		newSenderArray = append(newSenderArray, walletIdSenderArray[sendersNew[i]])
		newRecipientArray = append(newRecipientArray, walletIdSenderArray[recipientsNew[i]])
		newAmountArray = append(newAmountArray, newMirrorAmount[i])
	}
	return newSenderArray, newRecipientArray, newAmountArray
}

// Вспомогательная функция, которая проверяет, содержится ли число в массиве
func contains(arr []int, num int) bool {
	for _, n := range arr {
		if n == num {
			return true
		}
	}
	return false
}
