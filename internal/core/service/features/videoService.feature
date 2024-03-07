Feature: FindAll video service

    Scenario: when there is no video
      Given Admin post no video
      When admin run FindAll method
      Then video should be null

    Scenario: when there exist a video
      Given Admin post some video
      When admin run FindAll method
      Then video should be video

      

    
 