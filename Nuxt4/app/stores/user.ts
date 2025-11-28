import { FetchError } from 'ofetch'
export const useUserStore = defineStore(
  'user',
  () => {
    const id = ref<number | null>(null)
    const name = ref<string | null>(null)
    const icon = ref<string | null>(null)
    const coin = ref<number | null>(null)
    const date = ref<Date | null>(null)
    const paidSubscription = ref<boolean | null>(null)
    const dateSubstring = ref<Date | null>(null)

    const setNameIcon = (Name: string, Icon: string) => {
      name.value = Name
      icon.value = Icon
    }

    const getMySettings = async () => {
      try {
        const response = await useApiFetch<{
          id: number
          name: string
          icon: string | null
          coin: number
          dateOfRegistration: string
          paidSubscription: boolean
          dateOfPaidSubscription: string | null
        }>('/api/settings', {
          method: 'GET',
        })

        id.value = response.id
        name.value = response.name
        icon.value = response.icon
        coin.value = response.coin
        date.value = new Date(response.dateOfRegistration)
        paidSubscription.value = response.paidSubscription
        if(response.dateOfPaidSubscription)
          dateSubstring.value = new Date(response.dateOfPaidSubscription)
      } catch (error: unknown) {
         if(error instanceof FetchError) logout()
      }
    }
    // const setPaidSubscription = (Paid: boolean, Date: Date) => {
    //   paidSubscription.value = Paid
    //   dateSubstring.value = Date
    // }

    const logout = () => {
      id.value = null
      name.value = null
      icon.value = null
      coin.value = null
      date.value = null
      paidSubscription.value = null
      dateSubstring.value = null
    }

    const isLogged = computed(() => {return !!id.value})

    return { id, name, icon, coin, date, paidSubscription, dateSubstring, getMySettings, setNameIcon, logout, isLogged }
  },
  {
    persist: {
      pick: ['id', 'name', 'icon', 'coin'],
    },
  },
)
