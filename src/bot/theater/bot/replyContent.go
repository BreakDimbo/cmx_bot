package bot

var replySlice []string
var iteraSlice []string

func init() {
	replySlice = []string{
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
		"连10秒都用不了吧",
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

	iteraSlice = []string{
		"不是嘿客是黑客吧常考",
		"话说，差不多该放弃了吧。别对500日元买来的处理品软件抱期待啊",
		"那家伙根本就没有是现实还是游戏的想法吧",
		"那得另当别论，因为那些女孩都是我的老婆",
		"但是啊，是个很有意思的主题吧。假如我们是屏幕里的存在的话，有能够确认的方法吗",
		"即答啊",
		"又开始了厨二病。乙！辛苦了",
		"这是设定的吧",
		"今天的「只有不想被你这么说」原来是这里啊",
		"说起来，报道的直升机好像到现场了",
		"早上冈伦不是生气地说「那个博士居然临阵脱逃了」的吗",
		"都说了就没有举行啊。不是挺好的吗，举行了的话冈伦大概就成了人造卫星的肉垫了",
		"又是设定吗，从那个斯坦因什么的开始就莫名其妙",
		"给",
		"是白给的东西吧",
		"果断拒绝",
		"因为今天很热啊",
		"再说别以为胖子就一定有力气哦",
		"计划？是啥",
		"电话微波炉（暂定）啊",
		"话说，微波炉投入实战是什么",
		"高中二年级的时候因为班级不同而几乎没怎么说过话，所以实际上只有两年吧",
		"那·是·不·可·能·的",
		"话说，还要再试试吗？已经准备好了",
		"很漂亮的转盘吧，和平时不同这是在倒转啊",
		"有是有…才怪",
		"没有啊",
		"还是这样呢，既没有变热也没有变冷",
		"又来了",
		"不要",
		"真由氏，「你的香蕉软软的呢」说说看",
		"啊..我要融化了…",
		"夏天从距离上来说到菲伊莉斯炭那里就是极限喽，要能在「MayQueen+Nyan2」开就好了呢",
		"如果不关系到学分，就绝对不来了",
		"凉快，复活了",
		"你还在说这事？",
		"反正围观者太多根本就看不到吧，不过反正在@channel可以看到Snake的直播",
		"哦，暂时中止作业。又是祭典？",
		"为啥",
		"不允许侵犯个人隐私绝对…",
		"上星期发来的那个？",
		"还分了3条让人超不爽",
		"你看，Mail：「收信日期：7月23日12点56分」",
		"Mail：「牧濑红莉栖」",
		"Mail：「被什么人给捅了」",
		"Mail：「一刀好像」",
		"2…28号，怎么了",
		"喂… 冈伦？",
		"露易丝酱的名台词出现了这是",
		"冈伦，妄想只对我们说就够了，中钵博士的招待会中止了吧",
		"讲师原来是牧濑红莉栖啊",
		"你还在说吗",
		"小冈伦又乱来",
		"什么？我现在和菲伊莉斯炭…",
		"什么书？",
		"喂…冈伦",
		"约翰·提托是…谁啊？",
		"不可原谅，绝对的",
		"我想对赢家说了就等于败了",
		"今天的「轮不到你来说」原来是这里啊",
		"菲伊莉斯炭另当别论，无论是三次元还是二次元之魂都寄宿着",
		"哦，这是IBN5100吧",
		"梦幻的古董电脑",
		"差不多是一个月前吧，流传着它沉睡在秋叶原某处的传闻，听了那个才来的吧",
		"结果还是没找到，连疾风迅雷的Knight-Hart都出场了，但还是没找出来，最后结论就是从一开始就没有吧",
		"发售距今30年以上，当时是价格太高，谁都买不起电脑的时代，稀有度超高的例子",
		"怎么啦",
		"「世界有危险」来了！",
		"我开动了",
		"好吃好吃，这个真好吃",
		"就是说变成了不是香蕉的什么东西吗",
		"那岂不是糟了？",
		"发生了什么呢",
		"根据是？",
		"毫无根据…好了，连接完毕",
		"这…不是真由氏的吗",
		"所以说那就是真由氏的…",
		"一根就够了吧常考",
		"那个（暂定）太麻烦了差不多该去掉了",
		"想改变的只有冈伦而已",
		"就是这样啊",
		"好结束了",
		"嗯？哦？————！",
		"消失了，香蕉…",
		"自演，乙！",
		"是冈伦藏起来的吧",
		"什…什么？",
		"那怎么可能",
		"但…但是啊，这不是和把完全连在一起了吗",
		"大…大概",
		"是这么一回事吗",
		"桶子，你难道把我们出卖了吗",
		"被这个三次元女人的色相所迷惑",
		"冈伦的脑内设定，没什么特别的含义",
		"没那种事",
		"哈~~",
		"牧濑氏，牧濑氏",
		"「谁会吃变态的香蕉啊」能请你再说一遍么",
		"可以的话用很不甘心的表情",
		"嗯",
		"没啦～～",
		"但是牧濑氏的话也许能够解释这种古怪的机能哦",
		"再说我想只靠我们是绝对不可能弄清的",
		"冈伦真狭小，做人的器量太狭小了",
		"这真让人着迷让人崇拜",
		"变态同志的视线杀人战，够燃",
		"那个…就是说从微波炉里放出像闪电一样的电流",
		"是昨天中午吧，小冈伦看到人工卫星坠落的新闻跑出去后，我把自己的手机接在这个东西上",
		"牧濑氏被捅了的短信？那不是一周前发的么",
		"因为在做运行测试，是正在使用反转功能吧",
		"写什么呢？",
		"那么就折衷吧，「冈伦是HENTAI」",
		"那么我发喽",
		"换气换气",
		"地板上的洞怎么办，被布朗氏发现的话可不是能简单完事的",
		"哎...",
		"大受打击啊，通宵了居然还毫无成果，我对做的事感到后悔，为什么又做不出胶化香蕉了呢",
		"把牧濑氏叫回来比较好吧",
		"是怎么回事啊",
		"LHC，「Large  Hadron  Collider」SERN的基本粒子加速器",
		"被治愈了呢",
		"之前提到的家伙？",
		"现在@channel上就像是祭典一样呢",
		"意义不明啊",
		"都说了，是「欧洲原子核共同研究机构」的略称，总部在日内瓦郊外，正如其名主要是基本粒子物理学的研究，为此拥有着数座世界级的设施，低能反质子环，质子同步推进器，大型正负电子对撞器，然后最终BOSS是——全球最大的粒子加速器，Large  Hadron  Collider炭，虽然也有LHC根据用法也可以生成迷你黑洞的传闻",
		"咋了",
		"有是有，之前牧濑氏在演讲时说过吧，归根到底生成黑洞就是不可能的",
		"有是有…不，才怪，SERN官方都否定了",
		"你究竟在和谁战斗啊",
		"去菲伊莉斯炭那里…",
		"不是嘿客是黑客啊",
		"不，完全不明白",
		"冈伦，你在搞笑吧",
		"我可不知道会发生什么哦？",
		"就差一点了，只要弄到SQL的一览表找出密码就容易了…",
		"所以说不是这里",
		"不准说嘿客！",
		"来了！————————",
		"来了来了来了来了",
		"认输吧，全都展现在我的面前吧",
		"来了，ID到手！",
		"任务完成",
		"虽说成功入侵了但不是管理员ID的话，能看的范围也是有限制的",
		"只能看到胸部的感觉",
		"还只是出来了几封邮件而已",
		"根本就不用那么麻烦",
		"翻译出来就是：嗨波尔",
		"是嘛？",
		"Mail：「那个，LHC的状态良好，这家伙虽然像猫一样反复无常，最近一个月的状态却惊人得好」",
		"不知道…",
		"虽然没有时间机器这个单词，「Z  Program」这个单词在最近一个月被使用了100多次",
		"嗯…第137次 Z  Program 实验报告，因迷你黑洞生成任务已经完成，报告省略",
		"好像是的，官方明明宣称实验没有成功的呢",
		"实验结果",
		"「Human  is  Dead, mismatch」",
		"嗯",
		"实验结果「Human is Dead, mismatch」",
		"嗯",
		"我想是这样…大概",
		"那个，详情参阅「果冻人报告No.14」",
		"我怎么可能知道嘛，只是这边的服务器上有奇怪的数据库",
		"大概是某种程序代码",
		"啊哈哈哈哈！！搞得懂才怪，这种的绝对搞不懂，去死，渣滓！",
		"我本以为到这步就简单了可还是失败了",
		"就算你不说我也不行了",
		"多谢，给我来一罐",
		"是在说关东煮吗…",
		"你在说什么我不明白哦",
		"小冈伦你能赢吗？菲利斯酱炭可是全国水平哦",
		"作战？",
		"好快",
		"作战到底是什么啊",
		"怎么看都是冷笑话真是非常感谢",
		"什么",
		"该说是智慧性还是肮脏呢",
	}
}
