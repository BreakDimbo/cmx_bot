package bot

import (
	"bot/bredis"
	"bot/const"
	"bot/log"
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

var m *sync.Mutex
var replies map[string][]string

const FoodKey = "FoodKey"

func init() {
	replies = make(map[string][]string)
	replies[cons.Kurisu] = initKurisuScrip()
	replies[cons.Itaru] = initItaruScript()
}

func GetRandomReply(name string) string {
	m.Lock()
	defer m.Unlock()
	rs, ok := replies[name]
	if !ok {
		log.SLogger.Errorf("%s key not found.", name)
		return ""
	}
	index := rand.Intn(len(rs))
	return rs[index]
}

// AddReply add reply to replies map with lock
func AddReply(name string, reply string) error {
	m.Lock()
	defer m.Unlock()
	rs, ok := replies[name]
	if !ok {
		return fmt.Errorf("%s not exists error", name)
	}
	replies[name] = append(rs, reply)
	return nil
}

func initKurisuScrip() []string {
	return []string{
		"喂，我说你",
		"那是我的台词",
		"你知道得还挺多呢。你是哪个大学研究室里的人…",
		"你在跟谁说话啊",
		"是吗，自言自语么",
		"你想被警察抓走吗",
		"你是笨蛋吗，想死吗",
		"你是…",
		"あれ？",
		"要挑战看看吗，凤凰院先生？",
		"是嘛，那就请去找到异物质吧，凤凰院先生",
		"在做很有意思的实验呢",
		"冈部伦太郎…",
		"其实我是为了确认这个才来的",
		"但比起那个…",
		"里面黏糊糊的，味道呢…",
		"我才不要，谁会吃变态的香蕉啊",
		"突然来戳别人，还想连身体都摸",
		"看来你们两个都是变态的样子呢",
		"OK，说实话性骚扰已经到想马上去告发的地步了，但现在我不追究",
		"就是说你要邀请我成为研究所的成员么",
		"但是，我预定8月要回美国去",
		"还有么",
		"真是的，让人甲肾上腺素分泌过剩呢",
		"好好叫人名字，凤凰院变态凶真",
		"闭嘴变态！",
		"真没办法呢",
		"我知道了，我接受条件",
		"嗯",
		"是吗",
		"嗯，请多关照",
		"为什么刚才不说",
		"桥田先生，Good Job",
		"真的吗",
		"那怎么可能",
		"干嘛啊",
		"什…等等",
		"真是的，肩膀上的衣服都乱了嘛",
		"为什么要瞪着我啊",
		"要找人出气的话劝你还是放弃吧",
		"我可不想被你说",
		"比起这个，要说的话是什么啊",
		"为什么我就跟这种事扯上关系了呢",
		"真想把那时输给好奇心的自己揍一顿",
		"怎么听都是阴谋论，真是谢…",
		"没…没什么",
		"真的没什么",
		"要是再追究我要打人了",
		"你以为我会随便跟…",
		"你以为只要转一下就会发生什么吗",
		"那肯定是有什么地方搞错了",
		"我是不会参与那种东西的",
		"我说过了吧，理论上是不可能的",
		"不要",
		"我说了不要了",
		"抱歉，有点冲动了",
		"你在装什么帅啊",
		"就是这才让人生气",
		"我会离开的，再见",
		"好好，再见",
		"不要拉那里",
		"像傻瓜一样，真是白听了",
		"难道说是真的吗",
		"刚才不是才见过面吗",
		"是谁的错啊，谁的",
		"你现在在哪",
		"只是问下你在哪里",
		"好好，行了就算是助手吧",
		"那么你是？",
		"多大了",
		"比我小1岁啊",
		"什么意思",
		"又是那个吗",
		"果断拒绝",
		"为什么非得和你一起落到这个地步啊",
		"我说啊",
		"看得起我真是谢谢，但是不行",
		"要求并排走",
		"给我过来，要摔倒了",
	}
}

func initItaruScript() []string {
	itaruSlice := []string{
		"我不管，今晚吃炭烧鸡！",
		"不是嘿客是黑客吧常考",
		"那些女孩都是我的老婆",
		"又开始了厨二病。乙！辛苦了",
		"这是设定的吧",
		"今天的「只有不想被你这么说」原来是这里啊",
		"又是设定吗，从那个斯坦因什么的开始就莫名其妙",
		"果断拒绝",
		"再说别以为胖子就一定有力气哦",
		"那·是·不·可·能·的",
		"话说，还要再试试吗？已经准备好了",
		"又来了",
		"不要",
		"「你的香蕉软软的呢」说说看",
		"啊..我要融化了…",
		"什么？我现在和菲伊莉斯炭…",
		"不可原谅，绝对的",
		"我开动了",
		"好吃好吃，这个真好吃",
		"一根就够了吧常考",
		"嗯？哦？————！",
		"自演，乙！",
		"被这个三次元女人的色相所迷惑",
		"哈~~",
		"「谁会吃变态的香蕉啊」能请你再说一遍么，可以的话用很不甘心的表情",
		"没啦~~",
		"这真让人着迷让人崇拜",
		"变态同志的视线杀人战，够燃",
		"被治愈了呢",
		"意义不明啊",
		"咋了",
		"你究竟在和谁战斗啊",
		"去菲伊莉斯炭那里…",
		"不，完全不明白",
		"我可不知道会发生什么哦？",
		"不准说嘿客！",
		"来了！————————",
		"来了来了来了来了",
		"认输吧，全都展现在我的面前吧",
		"任务完成",
		"只能看到胸部的感觉",
		"啊哈哈哈哈！！搞得懂才怪，这种的绝对搞不懂，去死，渣滓！",
		"你在说什么我不明白哦",
		"怎么看都是冷笑话真是非常感谢",
	}

	keyPattern := FoodKey + "*"
	keys, err := bredis.Client.Keys(keyPattern).Result()
	if err != nil {
		log.SLogger.Errorf("get food key from redis error: %v", err)
	}

	var foods []string
	for _, key := range keys {
		keySlice := strings.Split(key, ":")
		food := fmt.Sprintf("诶嘿嘿，%s 怎么样？", keySlice[1])
		log.SLogger.Infof("food initia: %s", food)

		foods = append(foods, food)
	}

	itaruSlice = append(itaruSlice, foods...)
	return itaruSlice
}
