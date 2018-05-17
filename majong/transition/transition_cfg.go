package transition

// TODO 改成使用工具从 yaml 配置文件中生成此文件
var transitionCfg = `
- 
  game_id: 1
  # 状态表
  states:
    - 
      # 初始化状态
      state: state_init
      # 转换表
      transition:
        - 
        # 开始游戏事件 转移到 洗牌状态 
          events: 
            - event_start_game
          next_state: state_xipai

    - 
      # 洗牌状态
      state: state_xipai
      # 转换表
      transition:
        - 
        # 洗牌完成事件 转移到 发牌状态 
          events: 
            - event_xipai_finish
          next_state: state_fapai
    - 
      # 发牌状态
      state: state_fapai
      # 转换表
      transition:
        - 
        # 发牌完成事件 转移到 换三张状态
          events: 
            - event_fapai_finish
          next_state: state_huansanzhang
    - 
      # 换三张状态
      state: state_huansanzhang
      # 转换表
      transition:
        - 
          # 玩家换三张请求事件 转移到 换三张状态
          events: 
            - event_huansanzhang_request
          next_state: state_huansanzhang
        - 
          # 玩家换三张请求事件 转移到 定缺状态
          events: 
            - event_huansanzhang_request
          next_state: state_dingque
    - 
      # 定缺状态
      state: state_dingque
      # 转换表
      transition:
        - 
          # 玩家定缺请求事件 转移到 定缺状态
          events: 
            - event_dingque_request
          next_state: state_dingque
        - 
          # 玩家定缺请求事件 转移到 自询状态
          events: 
            - event_dingque_request
          next_state: state_zixun
    - 
      # 出牌状态
      state: state_chupai
      # 转换表
      transition:
        - 
          # 出牌完成事件 转移到 出牌问询状态
          events: 
            - event_chupai_finish
          next_state: state_chupaiwenxun
        - 
          # 出牌完成事件 转移到 摸牌状态
          events: 
            - event_chupai_finish
          next_state: state_mopai
    -
      # 出牌问询状态
      state: state_chupaiwenxun 
      # 转换表
      transition:
        -
        # 碰杠胡弃事件 转移到 出牌问询状态
          events:
            - event_peng_request
            - event_gang_request
            - event_hu_request
            - event_qi_request
          next_state: state_chupaiwenxun
        - 
        # 碰弃事件 转移到 碰状态
          events: 
            - event_peng_request
            - event_qi_request 
          next_state: state_peng 
        - 
        # 杠弃事件 转移到 杠状态
          events: 
            - event_gang_request
            - event_qi_request 
          next_state: state_gang 
        - 
        # 碰杠胡弃事件 转移到 胡状态
          events: 
            - event_peng_request
            - event_gang_request
            - event_hu_request
            - event_qi_request
          next_state: state_hu 
        - 
        # 弃事件 转移到 摸牌状态
          events: 
            - event_qi_request
          next_state: state_mopai
    - 
      # 暗杠状态
      state: state_angang
      # 转换表
      transition:
        - 
          # 暗杠完成事件 转移到 摸牌状态
          events: 
            - event_angang_finish
          next_state: state_mopai
    - 
      # 自摸状态
      state: state_zimo
      # 转换表
      transition:
        - 
          # 自摸完成事件 转移到 摸牌状态
          events: 
            - event_zimo_finish
          next_state: state_mopai
    - 
      # 碰状态
      state: state_peng
      # 转换表
      transition:
        - 
          # 玩家出牌事件 转移到 出牌状态
          events: 
            - event_chupai_request
          next_state: state_chupai 
    - 
      # 杠状态
      state: state_gang
      # 转换表
      transition:
        - 
          # 杠完成 转移到 摸牌状态
          events: 
            - event_gang_finish
          next_state: state_mopai  
    - 
      # 胡状态
      state: state_hu 
      # 转换表
      transition:
        - 
          # 胡完成事件 转移到 摸牌状态
          events: 
            - event_hu_finish
          next_state: state_mopai  
    - 
      # 摸牌状态
      state: state_mopai 
      # 转换表
      transition:
        - 
          # 摸牌完成事件 转移到 自询
          events: 
            - event_mopai_finish
          next_state: state_zixun   
        - 
          # 摸牌完成事件 转移到 结束
          events: 
            - event_mopai_finish
          next_state: state_gameover   
    - 
      # 自询状态
      state: state_zixun 
      # 转换表
      transition:
        - 
          # 玩家自摸请求事件 转移到 自摸状态
          events: 
            - event_hu_request
          next_state: state_zimo    
        - 
          # 玩家杠请求事件 转移到 暗杠状态
          events: 
            - event_gang_request
          next_state: state_angang   
        - 
          # 玩家出牌请求事件 转移到 出牌
          events: 
            - event_chupai_request
          next_state: state_chupai    
        - 
          # 玩家杠请求事件 转移到 补杠状态
          events: 
            - event_gang_request
          next_state: state_bugang     
        - 
          # 玩家杠请求事件 转移到 等待抢杠胡状态
          events: 
            - event_gang_request
          next_state: state_waitqiangganghu     
    - 
      # 补杠状态
      state: state_bugang 
      # 转换表
      transition:
        - 
          # 补杠完成 转移到 摸牌
          events: 
            - event_bugang_finish
          next_state: state_mopai  
    - 
      # 等待抢杠胡状态
      state: state_waitqiangganghu 
      # 转换表
      transition:
        - 
          # 玩家请求抢杠胡事件 或者 放弃抢杠胡请求事件 转移到 等待抢杠胡
          events: 
            - event_hu_request
            - event_qi_request
          next_state: state_waitqiangganghu   
        - 
          # 玩家请求抢杠胡事件  或者 放弃抢杠胡请求事件 转移到 抢杠胡
          events: 
            - event_hu_request
            - event_qi_request
          next_state: state_qiangganghu   
        - 
          # 玩家放弃抢杠胡请求事件 转移到 补杠
          events: 
            - event_qi_request
          next_state: state_bugang    
    - 
      # 抢杠胡状态
      state: state_qiangganghu  
      # 转换表
      transition:
        - 
          # 抢杠胡完成 转移到 摸牌
          events: 
            - event_qiangganghu_finish
          next_state: state_mopai  
`
